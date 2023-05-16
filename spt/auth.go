package spt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_auth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/aditya-K2/utils"
	"github.com/zmb3/spotify/v2"
)

var (
	redirectURI = "http://localhost:8080/callback"
	scopes      = []string{
		_auth.ScopeUserLibraryRead,
		_auth.ScopePlaylistModifyPublic,
		_auth.ScopePlaylistModifyPrivate,
		_auth.ScopePlaylistReadPrivate,
		_auth.ScopePlaylistReadCollaborative,
		_auth.ScopeUserReadPlaybackState,
		_auth.ScopeUserModifyPlaybackState,
		_auth.ScopeUserLibraryModify,
		_auth.ScopeUserLibraryRead,
		_auth.ScopeUserReadPrivate,
		_auth.ScopeUserFollowRead,
		_auth.ScopeUserReadCurrentlyPlaying,
		_auth.ScopeUserModifyPlaybackState,
		_auth.ScopeUserReadRecentlyPlayed,
		_auth.ScopeUserTopRead,
		_auth.ScopeStreaming,
	}
	auth = _auth.New(
		_auth.WithRedirectURL(redirectURI),
		_auth.WithScopes(scopes...))
	ch                           = make(chan *payload)
	state                        = "__GSPT_AUTH__"
	userConfigDir, userConfigErr = os.UserConfigDir()
	gsptDir                      = filepath.Join(userConfigDir, "/gspt")
	tokenPath                    = filepath.Join(gsptDir, "/oauthtoken")
)

type payload struct {
	Token *oauth2.Token
	Err   error
}

func InitClient() error {
	clientID := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")
	if clientID == "" || clientSecret == "" {
		return errors.New("SPOTIFY_ID and/or SPOTIFY_SECRET are missing. Please make sure you have set the SPOTIFY_ID and SPOTIFY_SECRET environment variables")
	}

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

		fmt.Println("Please log in to Spotify by visiting the following page in your browser: ")
		fmt.Println(url)

		// wait for auth to complete
		payload := <-ch

		if payload.Err != nil {
			return payload.Err
		}

		token = payload.Token
	}

	Client = spotify.New(auth.Client(context.Background(), token))
	return nil
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
		if !utils.FileExists(gsptDir) {
			if derr := os.Mkdir(gsptDir, 0777); derr != nil {
				ch <- &payload{nil, derr}
			}
		}
		if werr := os.WriteFile(tokenPath, val, 0777); werr != nil {
			ch <- &payload{nil, werr}
		}
	}

	ch <- &payload{tok, nil}
}
