package esi

type SkillList struct {
	Skills []Skill `json:"skills"`
}

type Skill struct {
	SkillID      int64 `json:"skill_id"`
	TrainedLevel int64 `json:"trained_skill_level"`
}
