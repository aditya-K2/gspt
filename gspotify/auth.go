package gspotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	sptauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/aditya-K2/gspot/utils"
	"github.com/zmb3/spotify/v2"
)

var (
	redirectURI = "http://localhost:8080/callback"
	scopes      = []string{
		sptauth.ScopeUserLibraryRead,
		sptauth.ScopePlaylistModifyPublic,
		sptauth.ScopePlaylistModifyPrivate,
		sptauth.ScopePlaylistReadCollaborative,
		sptauth.ScopeUserLibraryModify,
		sptauth.ScopeUserLibraryRead,
		sptauth.ScopeUserReadPrivate,
		sptauth.ScopeUserReadCurrentlyPlaying,
		sptauth.ScopeUserModifyPlaybackState,
		sptauth.ScopeUserReadRecentlyPlayed,
		sptauth.ScopeUserTopRead,
		sptauth.ScopeStreaming,
	}
	auth = sptauth.New(
		sptauth.WithRedirectURL(redirectURI),
		sptauth.WithScopes(scopes...))
	ch                           = make(chan *payload)
	state                        = "__GSPOT_AUTH__"
	userConfigDir, userConfigErr = os.UserConfigDir()
	gspotDir                     = filepath.Join(userConfigDir, "/gspot")
	tokenPath                    = filepath.Join(gspotDir, "/oauthtoken")
)

type payload struct {
	Token *oauth2.Token
	Err   error
}

func NewClient() (*spotify.Client, error) {
	token := &oauth2.Token{}

	// shouldn't be nil if the file doesn't exist.
	tokenErr := errors.New("")

	if utils.FileExists(tokenPath) {
		var content []byte
		content, tokenErr = os.ReadFile(tokenPath)
		tokenErr = json.Unmarshal(content, token)
	}

	if tokenErr != nil {
		http.HandleFunc("/callback", completeAuth)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Got request for:", r.URL.String())
		})

		go func() {
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				log.Fatal(err)
				ch <- &payload{nil, err}
			}
		}()
		url := auth.AuthURL(state)

		utils.Print("BLUE", "Please log in to Spotify by visiting the following page in your browser: ")
		utils.Print("GREEN", url)

		// wait for auth to complete
		payload := <-ch

		if payload.Err != nil {
			return nil, payload.Err
		}

		token = payload.Token
	}

	return spotify.New(auth.Client(context.Background(), token)), nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
		ch <- &payload{nil, err}
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		_s := fmt.Sprintf("State mismatch: %s != %s\n", st, state)
		log.Fatalf(_s)
		ch <- &payload{nil, errors.New(_s)}
	}

	if val, merr := json.Marshal(tok); merr != nil {
		ch <- &payload{nil, merr}
	} else {
		if !utils.FileExists(gspotDir) {
			if derr := os.Mkdir(gspotDir, 0777); derr != nil {
				ch <- &payload{nil, derr}
			}
		}
		if werr := os.WriteFile(tokenPath, val, 0777); werr != nil {
			ch <- &payload{nil, werr}
		}
	}

	ch <- &payload{tok, nil}
}
