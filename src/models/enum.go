package models

type Judge int

const (
	Codeforces Judge = iota + 1
	CodeforcesGym
	TopCoder
	AtCoder
	CsAcademy
	Codechef
	HackerRank
	HackerEarth
	Leetcode
	Others
)

// CodeForces, CodeForces::Gym, TopCoder, AtCoder, CS Academy, CodeChef, HackerRank, HackerEarth, Kick Start, LeetCode, Toph,
func (j Judge) String() *string {
	judges := [...]string{"CodeForces", "CodeForces::Gym", "TopCoder", "AtCoder", "CS Academy", "CodeChef", "HackerRank", "HackerEarth", "LeetCode", "Others"}

	if j < Codeforces || j > Others {
		return nil
	}

	return &judges[j-1]
}
