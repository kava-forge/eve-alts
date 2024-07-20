package characters

import (
	"context"
	"database/sql"

	"golang.org/x/oauth2"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"

	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/esi"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func RefreshCharacterData(ctx context.Context, deps dependencies, tok *oauth2.Token, charID int64) (repository.CharacterDBData, error) {
	logger := logging.With(deps.Logger(), keys.Component, "RefreshCharacterData")

	var data repository.CharacterDBData

	pubData, err := deps.ESIClient().GetCharacterPublicData(ctx, tok, charID)
	if err != nil {
		return data, errors.Wrap(err, "could not GetCharacterPublicData")
	}

	portraitData, err := deps.ESIClient().GetCharacterPortrait(ctx, tok, charID)
	if err != nil {
		return data, errors.Wrap(err, "could not GetCharacterPortrait")
	}

	corpData, err := deps.ESIClient().GetCorporationData(ctx, tok, pubData.CorporationID)
	if err != nil {
		return data, errors.Wrap(err, "could not GetCorporationData")
	}

	corpIcons, err := deps.ESIClient().GetCorporationIcons(ctx, tok, pubData.CorporationID)
	if err != nil {
		return data, errors.Wrap(err, "could not GetCorporationIcons")
	}

	skillList, err := deps.ESIClient().GetSkills(ctx, tok, charID)
	if err != nil {
		return data, errors.Wrap(err, "could not GetSkills")
	}
	skillIDMap := make(map[int64]bool, len(skillList.Skills))
	for _, skill := range skillList.Skills {
		skillIDMap[skill.SkillID] = true
	}

	seenSkills, err := deps.AppRepo().GetAllCharacterSkills(ctx, charID, nil)
	if err != nil && !errors.Is(err, database.ErrNoRows) {
		return data, errors.Wrap(err, "could not GetAllCharacterSkills")
	}

	toDelete := make([]int64, 0, 10)
	for _, skill := range seenSkills {
		if !skillIDMap[skill.SkillID] {
			toDelete = append(toDelete, skill.SkillID)
		}
	}

	var allianceData esi.AllianceData
	var allianceIcons esi.AllianceIcons
	if corpData.AllianceID != 0 {
		allianceData, err = deps.ESIClient().GetAllianceData(ctx, tok, corpData.AllianceID)
		if err != nil {
			return data, errors.Wrap(err, "could not GetAllianceData")
		}

		allianceIcons, err = deps.ESIClient().GetAllianceIcons(ctx, tok, corpData.AllianceID)
		if err != nil {
			return data, errors.Wrap(err, "could not GetAllianceIcons")
		}
	}

	var dbAlliance repository.Alliance
	var dbCorporation repository.Corporation
	var dbChar repository.Character
	var dbSkills []repository.CharacterSkill
	// var dbTok repository.Token
	if err := database.TransactWithRetries(ctx, deps.Telemetry(), logger, deps.DB(), &sql.TxOptions{}, func(ctx context.Context, tx database.Tx) error {
		var err error

		if dbSkills == nil {
			dbSkills = make([]repository.CharacterSkill, 0, len(skillList.Skills))
		} else {
			dbSkills = dbSkills[:0]
		}

		if corpData.AllianceID != 0 {
			if dbAlliance, err = deps.AppRepo().UpsertAlliance(ctx, corpData.AllianceID, allianceData.Name, allianceData.Ticker, allianceIcons.Small, tx); err != nil {
				return errors.Wrap(err, "could not UpsertAlliance")
			}
		}

		if dbCorporation, err = deps.AppRepo().UpsertCorporation(ctx, pubData.CorporationID, corpData.Name, corpData.Ticker, corpIcons.Small, corpData.AllianceID, tx); err != nil {
			return errors.Wrap(err, "could not UpsertCorporation")
		}

		if dbChar, err = deps.AppRepo().UpsertCharacter(ctx, charID, pubData.Name, portraitData.Medium, pubData.CorporationID, tx); err != nil {
			return errors.Wrap(err, "could not UpsertCharacter")
		}

		if _, err = deps.AppRepo().UpsertToken(ctx, dbChar.ID, tok.AccessToken, tok.RefreshToken, tok.TokenType, tok.Expiry, tx); err != nil {
			return errors.Wrap(err, "could not UpsertToken")
		}

		for _, skill := range skillList.Skills {
			dbSkill, err := deps.AppRepo().UpsertCharacterSkill(ctx, dbChar.ID, skill.SkillID, skill.TrainedLevel, tx)
			if err != nil {
				return errors.Wrap(err, "could not UpsertCharacterSkill")
			}
			dbSkills = append(dbSkills, dbSkill)
		}

		if len(toDelete) > 0 {
			if err := deps.AppRepo().DeleteCharacterSkills(ctx, dbChar.ID, toDelete, tx); err != nil {
				return errors.Wrap(err, "could not DeleteCharacterSkills")
			}
		}

		return nil
	}); err != nil {
		return data, errors.Wrap(err, "could not save character data")
	}

	data.Character = dbChar
	data.Corporation = dbCorporation
	data.Alliance = dbAlliance
	data.Skills = dbSkills

	return data, nil
}
