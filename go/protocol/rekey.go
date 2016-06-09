// Auto-generated by avdl-compiler v1.3.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/rekey.avdl

package keybase1

import (
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	context "golang.org/x/net/context"
)

type TLFID string
type TLF struct {
	Tlfid     TLFID    `codec:"tlfid" json:"tlfid"`
	Name      string   `codec:"name" json:"name"`
	Writers   []string `codec:"writers" json:"writers"`
	Readers   []string `codec:"readers" json:"readers"`
	IsPrivate bool     `codec:"isPrivate" json:"isPrivate"`
}

type ProblemTLF struct {
	Tlf       TLF   `codec:"tlf" json:"tlf"`
	Score     int   `codec:"score" json:"score"`
	Solutions []KID `codec:"solutions" json:"solutions"`
}

// ProblemSet is for a particular (user,kid) that initiated a rekey problem.
// This problem consists of one or more problem TLFs, which are individually scored
// and have attendant solutions --- devices that if they came online can rekey and
// solve the ProblemTLF.
type ProblemSet struct {
	User User         `codec:"user" json:"user"`
	Kid  KID          `codec:"kid" json:"kid"`
	Tlfs []ProblemTLF `codec:"tlfs" json:"tlfs"`
}

type ProblemSetDevices struct {
	ProblemSet ProblemSet `codec:"problemSet" json:"problemSet"`
	Devices    []Device   `codec:"devices" json:"devices"`
}

type Outcome int

const (
	Outcome_NONE    Outcome = 0
	Outcome_FIXED   Outcome = 1
	Outcome_IGNORED Outcome = 2
)

type ShowPendingRekeyStatusArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type GetPendingRekeyStatusArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type DebugShowRekeyStatusArg struct {
	SessionID int     `codec:"sessionID" json:"sessionID"`
	Tlfs      []TLFID `codec:"tlfs" json:"tlfs"`
	User      *UID    `codec:"user,omitempty" json:"user,omitempty"`
	Kid       *KID    `codec:"kid,omitempty" json:"kid,omitempty"`
}

type RekeyStatusFinishArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type RekeyInterface interface {
	// ShowPendingRekeyStatus shows either pending gregor-initiated rekey harassments
	// or nothing if none were pending.
	ShowPendingRekeyStatus(context.Context, int) error
	// GetPendingRekeyStatus returns the pending ProblemSetDevices.
	GetPendingRekeyStatus(context.Context, int) (ProblemSetDevices, error)
	// ShowRekeyStatus is used by the CLI to kick off a "ShowRekeyStatus" window for the given user based on
	// the passed-in parameters. These are the parameters that are typically delivered via direct
	// gregor injection. Will be used primarily in debugging or in advanced command-line usage.
	DebugShowRekeyStatus(context.Context, DebugShowRekeyStatusArg) error
	// rekeyStatusFinish is called when work is completed on a given RekeyStatus window. The Outcome
	// can be Fixed or Ignored.
	RekeyStatusFinish(context.Context, int) (Outcome, error)
}

func RekeyProtocol(i RekeyInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "keybase.1.rekey",
		Methods: map[string]rpc.ServeHandlerDescription{
			"showPendingRekeyStatus": {
				MakeArg: func() interface{} {
					ret := make([]ShowPendingRekeyStatusArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]ShowPendingRekeyStatusArg)
					if !ok {
						err = rpc.NewTypeError((*[]ShowPendingRekeyStatusArg)(nil), args)
						return
					}
					err = i.ShowPendingRekeyStatus(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"getPendingRekeyStatus": {
				MakeArg: func() interface{} {
					ret := make([]GetPendingRekeyStatusArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]GetPendingRekeyStatusArg)
					if !ok {
						err = rpc.NewTypeError((*[]GetPendingRekeyStatusArg)(nil), args)
						return
					}
					ret, err = i.GetPendingRekeyStatus(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"debugShowRekeyStatus": {
				MakeArg: func() interface{} {
					ret := make([]DebugShowRekeyStatusArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]DebugShowRekeyStatusArg)
					if !ok {
						err = rpc.NewTypeError((*[]DebugShowRekeyStatusArg)(nil), args)
						return
					}
					err = i.DebugShowRekeyStatus(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"rekeyStatusFinish": {
				MakeArg: func() interface{} {
					ret := make([]RekeyStatusFinishArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]RekeyStatusFinishArg)
					if !ok {
						err = rpc.NewTypeError((*[]RekeyStatusFinishArg)(nil), args)
						return
					}
					ret, err = i.RekeyStatusFinish(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type RekeyClient struct {
	Cli rpc.GenericClient
}

// ShowPendingRekeyStatus shows either pending gregor-initiated rekey harassments
// or nothing if none were pending.
func (c RekeyClient) ShowPendingRekeyStatus(ctx context.Context, sessionID int) (err error) {
	__arg := ShowPendingRekeyStatusArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.rekey.showPendingRekeyStatus", []interface{}{__arg}, nil)
	return
}

// GetPendingRekeyStatus returns the pending ProblemSetDevices.
func (c RekeyClient) GetPendingRekeyStatus(ctx context.Context, sessionID int) (res ProblemSetDevices, err error) {
	__arg := GetPendingRekeyStatusArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.rekey.getPendingRekeyStatus", []interface{}{__arg}, &res)
	return
}

// ShowRekeyStatus is used by the CLI to kick off a "ShowRekeyStatus" window for the given user based on
// the passed-in parameters. These are the parameters that are typically delivered via direct
// gregor injection. Will be used primarily in debugging or in advanced command-line usage.
func (c RekeyClient) DebugShowRekeyStatus(ctx context.Context, __arg DebugShowRekeyStatusArg) (err error) {
	err = c.Cli.Call(ctx, "keybase.1.rekey.debugShowRekeyStatus", []interface{}{__arg}, nil)
	return
}

// rekeyStatusFinish is called when work is completed on a given RekeyStatus window. The Outcome
// can be Fixed or Ignored.
func (c RekeyClient) RekeyStatusFinish(ctx context.Context, sessionID int) (res Outcome, err error) {
	__arg := RekeyStatusFinishArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.rekey.rekeyStatusFinish", []interface{}{__arg}, &res)
	return
}
