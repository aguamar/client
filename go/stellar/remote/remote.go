package remote

import (
	"context"
	"errors"
	"fmt"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/keybase1"
	"github.com/keybase/client/go/stellar/bundle"
)

type shouldCreateRes struct {
	Status       libkb.AppStatus `json:"status"`
	ShouldCreate bool            `json:"shouldcreate"`
}

func (r *shouldCreateRes) GetAppStatus() *libkb.AppStatus {
	return &r.Status
}

// ShouldCreate asks the server whether to create this user's initial wallet.
func ShouldCreate(ctx context.Context, g *libkb.GlobalContext) (should bool, err error) {
	defer g.CTraceTimed(ctx, "Stellar.ShouldCreate", func() error { return err })()
	arg := libkb.NewAPIArgWithNetContext(ctx, "stellar/shouldcreate")
	arg.SessionType = libkb.APISessionTypeREQUIRED
	var apiRes shouldCreateRes
	err = g.API.GetDecode(arg, &apiRes)
	return apiRes.ShouldCreate, err
}

// Post a bundle to the server
func PostWithChainlink(ctx context.Context, g *libkb.GlobalContext, clearBundle keybase1.StellarBundle) (err error) {
	defer g.CTraceTimed(ctx, "Stellar.Post", func() error { return err })()

	uid := g.ActiveDevice.UID()
	if uid.IsNil() {
		return libkb.NoUIDError{}
	}

	g.Log.CDebugf(ctx, "StellarPost load self")
	loadMeArg := libkb.NewLoadUserArg(g).
		WithNetContext(ctx).
		WithUID(uid).
		WithSelf(true).
		WithPublicKeyOptional()
	me, err := libkb.LoadUser(loadMeArg)
	if err != nil {
		return err
	}

	sigKey, err := g.ActiveDevice.SigningKey()
	if err != nil {
		return fmt.Errorf("signing key not found: (%v)", err)
	}
	pukring, err := g.GetPerUserKeyring()
	if err != nil {
		return err
	}
	err = pukring.Sync(ctx)
	if err != nil {
		return err
	}
	pukGen := pukring.CurrentGeneration()
	pukSeed, err := pukring.GetSeedByGeneration(ctx, pukGen)
	if err != nil {
		return err
	}

	err = clearBundle.CheckInvariants()
	if err != nil {
		return err
	}
	// Find the new primary account for the chain link.
	if len(clearBundle.Accounts) < 1 {
		return errors.New("stellar bundle has no accounts")
	}
	stellarAccount := clearBundle.Accounts[0]
	if len(stellarAccount.Signers) < 1 {
		return errors.New("stellar bundle has no signers")
	}
	if !stellarAccount.IsPrimary {
		return errors.New("initial stellar account is not primary")
	}
	g.Log.CDebugf(ctx, "StellarPost accountID:%v pukGen:%v", stellarAccount.AccountID, pukGen)
	boxed, err := bundle.Box(clearBundle, pukGen, pukSeed)
	if err != nil {
		return err
	}

	g.Log.CDebugf(ctx, "StellarPost make sigs")

	sig, err := libkb.WalletProofReverseSigned(me, stellarAccount.AccountID, stellarAccount.Signers[0], sigKey)
	if err != nil {
		return err
	}

	var sigsList []libkb.JSONPayload
	sigsList = append(sigsList, sig)

	payload := make(libkb.JSONPayload)
	payload["sigs"] = sigsList

	libkb.AddWalletServerArg(payload, boxed.EncB64, boxed.VisB64, int(boxed.FormatVersion))

	g.Log.CDebugf(ctx, "StellarPost post")
	_, err = g.API.PostJSON(libkb.APIArg{
		Endpoint:    "key/multi",
		SessionType: libkb.APISessionTypeREQUIRED,
		JSONPayload: payload,
	})
	if err != nil {
		return err
	}

	g.UserChanged(uid)
	return nil
}

type fetchRes struct {
	Status       libkb.AppStatus `json:"status"`
	EncryptedB64 string          `json:"encrypted"`
	VisibleB64   string          `json:"visible"`
}

func (r *fetchRes) GetAppStatus() *libkb.AppStatus {
	return &r.Status
}

// Fetch and unbox the latest bundle from the server.
func Fetch(ctx context.Context, g *libkb.GlobalContext) (res keybase1.StellarBundle, err error) {
	defer g.CTraceTimed(ctx, "Stellar.Fetch", func() error { return err })()
	arg := libkb.NewAPIArgWithNetContext(ctx, "stellar/bundle")
	arg.SessionType = libkb.APISessionTypeREQUIRED
	var apiRes fetchRes
	err = g.API.GetDecode(arg, &apiRes)
	if err != nil {
		return res, err
	}
	decodeRes, err := bundle.Decode(apiRes.EncryptedB64)
	if err != nil {
		return res, err
	}
	pukring, err := g.GetPerUserKeyring()
	if err != nil {
		return res, err
	}
	err = pukring.Sync(ctx)
	if err != nil {
		return res, err
	}
	puk, err := pukring.GetSeedByGeneration(ctx, decodeRes.Enc.Gen)
	if err != nil {
		return res, err
	}
	res, _, err = bundle.Unbox(decodeRes, apiRes.VisibleB64, puk)
	return res, err
}
