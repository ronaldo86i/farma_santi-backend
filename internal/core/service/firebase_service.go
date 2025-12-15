package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseClient Estructura que contiene la configuración de Firebase
type FirebaseClient struct {
	Ctx        context.Context
	App        *firebase.App
	AuthClient *auth.Client
}

const firebaseConfigFile = "./serviceAccountKey.json"

var (
	fbInstance *FirebaseClient
	once       sync.Once
)

// GetFirebaseClient devuelve una única instancia de firebaseClient
func GetFirebaseClient() *FirebaseClient {
	once.Do(func() {
		ctx := context.Background()
		opt := option.WithCredentialsFile(firebaseConfigFile)

		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			log.Fatalf("Error al inicializar firebase: %v\n", err)
		}
		log.Println("Firebase inicializado")
		authClient, err := app.Auth(ctx)
		if err != nil {
			log.Fatalf("Error inicializando AuthClient: %v\n", err)
		}
		log.Println("AuthClient inicializado")
		fbInstance = &FirebaseClient{
			Ctx:        ctx,
			App:        app,
			AuthClient: authClient,
		}
	})
	return fbInstance
}

// ObtenerAuthToken obtiene y verifica el token de autorización
func (f *FirebaseClient) ObtenerAuthToken(token string) (*auth.Token, error) {
	tokenDecoded, err := f.AuthClient.VerifyIDToken(f.Ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token no válido: %s", err.Error())
	}
	return tokenDecoded, nil
}
