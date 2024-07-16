package esi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	stdhttp "net/http"
	"net/url"
	"time"

	"github.com/hashicorp/cap/jwt"
	"github.com/kava-forge/eve-alts/lib/deferutil"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/json"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/browser"
	"github.com/rs/xid"
	"golang.org/x/oauth2"

	"github.com/kava-forge/eve-alts/pkg/repository"
)

const (
	ClientID  = "5a58af6b66a34b45a8b827e34b81527f"
	AuthURL   = "https://login.eveonline.com/v2/oauth/authorize/"
	TokenURL  = "https://login.eveonline.com/v2/oauth/token" //nolint:gosec,gocritic,revive
	UserAgent = "eve-alts (eve@evogames.org)"
)

var ErrInvalidCallbackCode = errors.New("invalid callback code")

var BaseURL, _ = url.Parse("https://esi.evetech.net")

type dependencies interface {
	Logger() logging.Logger
	ESICallbackServer() *CallbackServer
}

func TokenFromRepository(in repository.Token) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  in.AccessToken,
		RefreshToken: in.RefreshToken,
		TokenType:    in.TokenType,
		Expiry:       in.Expiration,
	}
}

type client struct {
	deps             dependencies
	callbackHostport string
	callbackPath     string
	oauth2           *oauth2.Config
	jwks             *jwt.KeySet
}

type Client interface {
	Authenticate(ctx context.Context) (*oauth2.Token, error)
	ValidateToken(ctx context.Context, tok *oauth2.Token) (CharacterData, error)
	GetCharacterPublicData(ctx context.Context, tok *oauth2.Token, charID int64) (CharacterPublicData, error)
	GetCharacterPortrait(ctx context.Context, tok *oauth2.Token, charID int64) (CharacterPortait, error)
	GetCorporationData(ctx context.Context, tok *oauth2.Token, corpID int64) (CorporationData, error)
	GetCorporationIcons(ctx context.Context, tok *oauth2.Token, corpID int64) (CorporationIcons, error)
	GetAllianceData(ctx context.Context, tok *oauth2.Token, allianceID int64) (AllianceData, error)
	GetAllianceIcons(ctx context.Context, tok *oauth2.Token, allianceID int64) (AllianceIcons, error)
	GetSkills(ctx context.Context, tok *oauth2.Token, charID int64) (SkillList, error)
}

func NewClient(deps dependencies, redirect string) (Client, error) {
	conf := &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
		RedirectURL: redirect,
		Scopes: []string{
			"publicData",
			"esi-skills.read_skills.v1",
			"esi-skills.read_skillqueue.v1",
			"esi-clones.read_clones.v1",
			"esi-assets.read_assets.v1",
		},
	}

	jwks, err := jwt.NewJSONWebKeySet(context.Background(), "https://login.eveonline.com/oauth/jwks", "")
	if err != nil {
		return nil, errors.Wrap(err, "could not load eve jwks")
	}

	return &client{
		deps:   deps,
		oauth2: conf,
		jwks:   &jwks,
	}, nil
}

func (c *client) Authenticate(ctx context.Context) (*oauth2.Token, error) {
	state := xid.New().String()
	verifier := oauth2.GenerateVerifier()
	u := c.oauth2.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))

	codeChan := make(chan CodeState)
	c.deps.ESICallbackServer().Expect(state, codeChan)
	defer c.deps.ESICallbackServer().Remove(state)

	egCtx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	if err := browser.OpenURL(u); err != nil {
		return nil, errors.Wrap(err, "could not open browser")
	}

	var code CodeState
	select {
	case code = <-codeChan:
	case <-egCtx.Done():
	}

	if !code.Valid {
		return nil, ErrInvalidCallbackCode
	}

	token, err := c.oauth2.Exchange(ctx, code.Code, oauth2.VerifierOption(verifier))
	return token, errors.Wrap(err, "could not exchange code for token")
}

func (c *client) ValidateToken(ctx context.Context, tok *oauth2.Token) (data CharacterData, err error) {
	v, err := jwt.NewValidator(*c.jwks)
	if err != nil {
		return data, errors.Wrap(err, "could not instantiate validator")
	}

	d, err := v.Validate(ctx, tok.AccessToken, jwt.Expected{
		Issuer: "https://login.eveonline.com",
		Audiences: []string{
			c.oauth2.ClientID,
			"EVE Online",
		},
	})
	if err != nil {
		return data, errors.Wrap(err, "invalid token")
	}

	if err := mapstructure.Decode(d, &data); err != nil {
		return data, errors.Wrap(err, "could not parse claims")
	}

	if err := data.fillRealID(); err != nil {
		return data, errors.Wrap(err, "could not parse character id")
	}

	return data, nil
}

func (c *client) makeRequest(ctx context.Context, tok *oauth2.Token, req *stdhttp.Request, target interface{}) error {
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.oauth2.Client(ctx, tok).Do(req)
	if err != nil {
		return errors.Wrap(err, "could not do http request")
	}
	defer deferutil.CheckDeferLog(c.deps.Logger(), resp.Body.Close)

	if resp.StatusCode != stdhttp.StatusOK {
		return errors.WithDetails(errors.New("bad request response"), "uri", req.URL.String(), "code", resp.StatusCode)
	}

	buf := &bytes.Buffer{}
	tee := io.TeeReader(resp.Body, buf)

	if err := json.UnmarshalFromReader(tee, target); err != nil {
		return errors.Wrap(err, "could not unmarshal response")
	}

	level.Debug(c.deps.Logger()).Message("response body", "uri", req.URL.String(), "contents", buf.String())

	return nil
}

func (c *client) GetCharacterPublicData(ctx context.Context, tok *oauth2.Token, charID int64) (CharacterPublicData, error) {
	u, err := BaseURL.Parse(fmt.Sprintf("/latest/characters/%d/", charID))
	if err != nil {
		return CharacterPublicData{}, errors.Wrap(err, "could not parse public data url")
	}

	req, err := stdhttp.NewRequestWithContext(ctx, stdhttp.MethodGet, u.String(), stdhttp.NoBody)
	if err != nil {
		return CharacterPublicData{}, errors.Wrap(err, "could not form http request")
	}

	var respData CharacterPublicData
	if err := c.makeRequest(ctx, tok, req, &respData); err != nil {
		return CharacterPublicData{}, errors.Wrap(err, "could not unmarshal response")
	}

	return respData, nil
}

func (c *client) GetCharacterPortrait(ctx context.Context, tok *oauth2.Token, charID int64) (CharacterPortait, error) {
	u, err := BaseURL.Parse(fmt.Sprintf("/latest/characters/%d/portrait/", charID))
	if err != nil {
		return CharacterPortait{}, errors.Wrap(err, "could not parse public data url")
	}

	req, err := stdhttp.NewRequestWithContext(ctx, stdhttp.MethodGet, u.String(), stdhttp.NoBody)
	if err != nil {
		return CharacterPortait{}, errors.Wrap(err, "could not form http request")
	}

	var respData CharacterPortait
	if err := c.makeRequest(ctx, tok, req, &respData); err != nil {
		return CharacterPortait{}, errors.Wrap(err, "could not unmarshal response")
	}

	return respData, nil
}

func (c *client) GetCorporationData(ctx context.Context, tok *oauth2.Token, corpID int64) (CorporationData, error) {
	u, err := BaseURL.Parse(fmt.Sprintf("/latest/corporations/%d/", corpID))
	if err != nil {
		return CorporationData{}, errors.Wrap(err, "could not parse public data url")
	}

	req, err := stdhttp.NewRequestWithContext(ctx, stdhttp.MethodGet, u.String(), stdhttp.NoBody)
	if err != nil {
		return CorporationData{}, errors.Wrap(err, "could not form http request")
	}

	var respData CorporationData
	if err := c.makeRequest(ctx, tok, req, &respData); err != nil {
		return CorporationData{}, errors.Wrap(err, "could not unmarshal response")
	}

	return respData, nil
}

func (c *client) GetCorporationIcons(ctx context.Context, tok *oauth2.Token, corpID int64) (CorporationIcons, error) {
	u, err := BaseURL.Parse(fmt.Sprintf("/latest/corporations/%d/icons/", corpID))
	if err != nil {
		return CorporationIcons{}, errors.Wrap(err, "could not parse public data url")
	}

	req, err := stdhttp.NewRequestWithContext(ctx, stdhttp.MethodGet, u.String(), stdhttp.NoBody)
	if err != nil {
		return CorporationIcons{}, errors.Wrap(err, "could not form http request")
	}

	var respData CorporationIcons
	if err := c.makeRequest(ctx, tok, req, &respData); err != nil {
		return CorporationIcons{}, errors.Wrap(err, "could not unmarshal response")
	}

	return respData, nil
}

func (c *client) GetAllianceData(ctx context.Context, tok *oauth2.Token, allianceID int64) (AllianceData, error) {
	u, err := BaseURL.Parse(fmt.Sprintf("/latest/alliances/%d/", allianceID))
	if err != nil {
		return AllianceData{}, errors.Wrap(err, "could not parse public data url")
	}

	req, err := stdhttp.NewRequestWithContext(ctx, stdhttp.MethodGet, u.String(), stdhttp.NoBody)
	if err != nil {
		return AllianceData{}, errors.Wrap(err, "could not form http request")
	}

	var respData AllianceData
	if err := c.makeRequest(ctx, tok, req, &respData); err != nil {
		return AllianceData{}, errors.Wrap(err, "could not unmarshal response")
	}

	return respData, nil
}

func (c *client) GetAllianceIcons(ctx context.Context, tok *oauth2.Token, allianceID int64) (AllianceIcons, error) {
	u, err := BaseURL.Parse(fmt.Sprintf("/latest/alliances/%d/icons/", allianceID))
	if err != nil {
		return AllianceIcons{}, errors.Wrap(err, "could not parse public data url")
	}

	req, err := stdhttp.NewRequestWithContext(ctx, stdhttp.MethodGet, u.String(), stdhttp.NoBody)
	if err != nil {
		return AllianceIcons{}, errors.Wrap(err, "could not form http request")
	}

	var respData AllianceIcons
	if err := c.makeRequest(ctx, tok, req, &respData); err != nil {
		return AllianceIcons{}, errors.Wrap(err, "could not unmarshal response")
	}

	return respData, nil
}

func (c *client) GetSkills(ctx context.Context, tok *oauth2.Token, charID int64) (SkillList, error) {
	u, err := BaseURL.Parse(fmt.Sprintf("/latest/characters/%d/skills/", charID))
	if err != nil {
		return SkillList{}, errors.Wrap(err, "could not parse skills url")
	}

	req, err := stdhttp.NewRequestWithContext(ctx, stdhttp.MethodGet, u.String(), stdhttp.NoBody)
	if err != nil {
		return SkillList{}, errors.Wrap(err, "could not form http request")
	}

	var respData SkillList
	if err := c.makeRequest(ctx, tok, req, &respData); err != nil {
		return SkillList{}, errors.Wrap(err, "could not unmarshal response")
	}

	return respData, nil
}
