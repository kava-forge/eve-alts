package characters

import (
	"github.com/kava-forge/eve-alts/pkg/operators"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func CharacterMatchesTag(char *repository.CharacterDBData, tag *repository.TagDBData) (bool, []repository.TagSkill) {
	charSkills := make(map[int64]int64, len(char.Skills))
	for _, sk := range char.Skills {
		charSkills[sk.SkillID] = sk.SkillLevel
	}

	return skillsMatchTag(charSkills, tag)
}

func skillsMatchTag(charSkills map[int64]int64, tag *repository.TagDBData) (bool, []repository.TagSkill) {
	missing := make([]repository.TagSkill, 0, len(tag.Skills))

	for _, sk := range tag.Skills {
		lvl, ok := charSkills[sk.SkillID]
		if !ok || lvl < sk.SkillLevel {
			missing = append(missing, sk)
		}
	}

	return len(missing) == 0, missing
}

func CharacterMatchesRole(char *repository.CharacterDBData, role *repository.RoleDBData, tags []*repository.TagDBData) (bool, []repository.Tag) {
	charSkills := make(map[int64]int64, len(char.Skills))
	for _, sk := range char.Skills {
		charSkills[sk.SkillID] = sk.SkillLevel
	}

	tLookup := make(map[int64]*repository.TagDBData)
	for _, tdb := range tags {
		tLookup[tdb.Tag.ID] = tdb
	}

	mismatch := make([]repository.Tag, 0, len(role.Tags))

	for _, t := range role.Tags {
		tdb, ok := tLookup[t.ID]
		if !ok {
			continue
		}

		match, _ := skillsMatchTag(charSkills, tdb)
		switch role.Operator {
		case operators.OperatorAny:
			if match {
				return true, nil
			} else {
				mismatch = append(mismatch, t)
			}
		case operators.OperatorAll:
			if !match {
				mismatch = append(mismatch, t)
			}
		case operators.OperatorNone:
			if match {
				mismatch = append(mismatch, t)
			}
		}
	}

	return len(mismatch) == 0, mismatch
}
