package lexer

import (
	"fmt"

	"github.com/mperham/inspeqtor-pro/jobs/token"
)

type ActionTable [NumStates]ActionRow

type ActionRow struct {
	Accept token.Type
	Ignore string
}

func (this ActionRow) String() string {
	return fmt.Sprintf("Accept=%d, Ignore=%s", this.Accept, this.Ignore)
}

var ActTab = ActionTable{
	ActionRow{ // S0
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S1
		Accept: -1,
		Ignore: "!whitespace",
	},
	ActionRow{ // S2
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S3
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S4
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S5
		Accept: 6,
		Ignore: "",
	},
	ActionRow{ // S6
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S7
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S8
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S9
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S10
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S11
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S12
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S13
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S14
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S15
		Accept: -1,
		Ignore: "!comment",
	},
	ActionRow{ // S16
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S17
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S18
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S19
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S20
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S21
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S22
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S23
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S24
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S25
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S26
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S27
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S28
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S29
		Accept: 11,
		Ignore: "",
	},
	ActionRow{ // S30
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S31
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S32
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S33
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S34
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S35
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S36
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S37
		Accept: 12,
		Ignore: "",
	},
	ActionRow{ // S38
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S39
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S40
		Accept: 7,
		Ignore: "",
	},
	ActionRow{ // S41
		Accept: 3,
		Ignore: "",
	},
	ActionRow{ // S42
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S43
		Accept: 4,
		Ignore: "",
	},
	ActionRow{ // S44
		Accept: 2,
		Ignore: "",
	},
	ActionRow{ // S45
		Accept: 14,
		Ignore: "",
	},
	ActionRow{ // S46
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S47
		Accept: 8,
		Ignore: "",
	},
	ActionRow{ // S48
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S49
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S50
		Accept: 9,
		Ignore: "",
	},
	ActionRow{ // S51
		Accept: 13,
		Ignore: "",
	},
	ActionRow{ // S52
		Accept: 10,
		Ignore: "",
	},
}