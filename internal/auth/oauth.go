package auth

import (
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/battlenet"
)

// InitBattleNetOAuth inicijuoja Battle.net OAuth su teisingais parametrais
func InitBattleNetOAuth() {
	clientID := os.Getenv("BATTLE_CLIENT_ID")
	secret := os.Getenv("BATTLE_CLIENT_SECRET")
	redirectURL := "http://localhost:3000/auth/battlenet/callback"
	scope := "wow.profile"                   // Tik wow.profile – ne locale!
	region := os.Getenv("BATTLE_NET_REGION") // pvz. "eu"

	goth.UseProviders(
		battlenet.New(clientID, secret, redirectURL, scope, region),
	)
}
