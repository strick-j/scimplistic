package views

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/db"
	"github.com/strick-j/scimplistic/handlers"
	"github.com/strick-j/scimplistic/internal/types"
	"github.com/strick-j/scimplistic/utils"
	"golang.org/x/oauth2"
)

// Object is used by both ScimType functions to deliver instructions
// and payload.
// Safe Add Example:
// 		safeObject := Object{
//			Method: "POST"
//			Type2Resources: addSafeStruct
//      }
// Note: In the above example the addSafeStruct should be a types.Type2Resources
// struct at least the required information to add a safe (e.g. Name and Schema)
// Docs for Example: https://identity-developer.cyberark.com/docs/manage-containers-with-scim-endpoints#post
type Object struct {
	Method         string // "GET", "PUT", "POST", "DELETE"
	Type           string // e.g. users, groups, containers, etc...
	Id             string // e.g. user id, group id, safe DisplayName, etc...
	Type1Resources types.Type1Resources
	Type2Resources types.Type2Resources
}

var (
	Type1          types.ScimType1
	Type2          types.ScimType2
	Type1Resources types.Type1Resources
	Type2Resources types.Type2Resources
)

////////// Service Functions //////////////////////////////////////////////////////////////////
type Service struct {
	client *handlers.Client
}

type transport struct {
	token string
}

func New(authToken *oauth2.Token) *Service {
	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		fmt.Println(err)
	}

	t := transport{
		token: authToken.AccessToken,
	}

	return &Service{
		client: handlers.New(
			&http.Client{Transport: &t},
			handlers.Options{
				ApiURL:  fmt.Sprintf("https://%s/%s/%s", values.ScimURL, values.ApiEndpoint, values.ApiVersion),
				Verbose: values.ServiceLogging,
			},
		),
	}
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", "Bearer "+t.token)

	return http.DefaultTransport.RoundTrip(r)
}

////////// Scim Type1 Functions ///////////////////////////////////////////////////////////////////
// ScimType1Api should be used for the following:
//   - Users
//   - Groups
// For Users and Groups ScimType1Api either returns the larger struct types.ScimType1
// or in the case of a single id GET/POST/PUT returns the sub struct types.ScimType1Resources
//
// Delete does not return a response is succesful, a response is only returned on error
////////////////////////////////////////////////////////////////////////////////////////////////////
func (ob Object) ScimType1Api() (*types.ScimType1, *types.Type1Resources, error) {
	ctx := context.Background()
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "ScimType1Api"}).Info("Processing Request")

	// Retrieve Authentication Token from handlers
	authToken := handlers.OauthCredClient()

	// Create Service Client with Authentication Token
	// Service Client will be used to make subsequent requests
	s := New(authToken)

	switch {
	case ob.Method == "GET":
		if ob.Id == "" {
			res, err := s.GetType1Objects(ctx, ob)
			if err != nil {
				return nil, nil, err
			}
			return res, nil, nil
		} else if ob.Id != "" {
			res, err := s.GetType1Object(ctx, ob)
			if err != nil {
				return nil, nil, err
			}
			return res, nil, nil
		}
	case ob.Method == "POST":
		res, err := s.AddType1Object(ctx, ob)
		if err != nil {
			return nil, nil, err
		}
		return nil, res, nil
	case ob.Method == "PUT":
		res, err := s.UpdateType1Object(ctx, ob)
		if err != nil {
			return nil, nil, err
		}
		return nil, res, nil
	case ob.Method == "DELETE":
		if ob.Id == "" {
			return nil, nil, fmt.Errorf("user id required for DELETE")
		} else if ob.Id != "" {
			if err := s.DeleteType1Object(ctx, ob); err != nil {
				return nil, nil, err
			}
			// If no error return nothing
			return nil, nil, nil
		}
	default:
		return nil, nil, fmt.Errorf("no method provided")
	}

	return nil, nil, nil
}

func (s *Service) GetType1Objects(ctx context.Context, ob Object) (*types.ScimType1, error) {
	if err := s.client.Get(ctx, fmt.Sprintf("/%s", ob.Type), &Type1); err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return &Type1, nil
}

func (s *Service) GetType1Object(ctx context.Context, ob Object) (*types.ScimType1, error) {
	type1query := url.PathEscape("filter=displayName eq " + ob.Id)
	if err := s.client.Get(ctx, fmt.Sprintf("/%s?%s", ob.Type, type1query), &Type1); err != nil {
		return nil, fmt.Errorf("failed to get %s %s: %w", ob.Type, ob.Id, err)
	}
	return &Type1, nil
}

func (s *Service) AddType1Object(ctx context.Context, ob Object) (*types.Type1Resources, error) {
	if err := s.client.Post(ctx, fmt.Sprintf("/%s", ob.Type), ob.Type1Resources, &Type1); err != nil {
		return nil, fmt.Errorf("failed to add %s %s: %w", ob.Type, ob.Type1Resources.UserName, err)
	}
	db.AddAction(ob.Method, ob.Type, Type1Resources.ID, "Success")
	return &Type1Resources, nil
}

func (s *Service) UpdateType1Object(ctx context.Context, ob Object) (*types.Type1Resources, error) {
	if err := s.client.Put(ctx, fmt.Sprintf("/%s/%s", ob.Type, ob.Id), ob.Type1Resources, &Type1); err != nil {
		return nil, fmt.Errorf("failed to update %s %s: %w", ob.Type, ob.Id, err)
	}
	return &Type1Resources, nil
}

func (s *Service) DeleteType1Object(ctx context.Context, ob Object) error {
	if err := s.client.Delete(ctx, fmt.Sprintf("/%s/%s", ob.Type, ob.Id), nil); err != nil {
		db.AddAction(ob.Method, ob.Type, ob.Id, "Failure")
		return fmt.Errorf("failed to delete %s %s: %w", ob.Type, ob.Id, err)
	}
	db.AddAction(ob.Method, ob.Type, ob.Id, "Success")
	return nil
}

////////// Scim Type2 Functions ///////////////////////////////////////////////////////////////////
// ScimType2Api should be used for the following:
//   - Containers
//   - Privileged Data
//   - Schemas
//   - ResourceTypes
// For Containers and Privileged Data ScimType2Api either returns the larger struct types.ScimType2
// or in the case of a single id GET/POST/PUT returns the sub struct types.ScimType2Resources
//
// Delete does not return a response is succesful, a response is only returned on error
////////////////////////////////////////////////////////////////////////////////////////////////////
func (ob Object) ScimType2Api() (*types.ScimType2, *types.Type2Resources, error) {
	ctx := context.Background()
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "ScimType2Api"}).Info("Processing Request")

	// Retrieve Authentication Token from handlers
	authToken := handlers.OauthCredClient()

	// Create Service Client with Authentication Token
	// Service Client will be used to make subsequent requests
	s := New(authToken)

	switch {
	case ob.Method == "GET":
		if ob.Id == "" {
			res, err := s.GetType2Objects(ctx, ob)
			if err != nil {
				return nil, nil, err
			}
			return res, nil, nil
		} else if ob.Id != "" {
			res, err := s.GetType2Object(ctx, ob)
			if err != nil {
				return nil, nil, err
			}
			return nil, res, nil
		}
	case ob.Method == "POST":
		res, err := s.AddType2Object(ctx, ob)
		if err != nil {
			return nil, nil, err
		}
		return nil, res, nil
	case ob.Method == "PUT":
		res, err := s.UpdateType2Object(ctx, ob)
		if err != nil {
			return nil, nil, err
		}
		return nil, res, nil
	case ob.Method == "DELETE":
		if ob.Id == "" {
			return nil, nil, fmt.Errorf("object ID required for DELETE")
		} else if ob.Id != "" {
			if err := s.DeleteType2Object(ctx, ob); err != nil {
				return nil, nil, err
			}
			// If no error return nothing
			return nil, nil, nil
		}
	default:
		return nil, nil, fmt.Errorf("no method provided")
	}

	return nil, nil, nil
}

func (s *Service) GetType2Objects(ctx context.Context, ob Object) (*types.ScimType2, error) {
	if err := s.client.Get(ctx, fmt.Sprintf("/%s", ob.Type), &Type2); err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return &Type2, nil
}

func (s *Service) GetType2Object(ctx context.Context, ob Object) (*types.Type2Resources, error) {
	if err := s.client.Get(ctx, fmt.Sprintf("/%s/%s", ob.Type, ob.Id), &Type2); err != nil {
		return nil, fmt.Errorf("failed to get %s %s: %w", ob.Type, ob.Id, err)
	}
	return &Type2Resources, nil
}

func (s *Service) AddType2Object(ctx context.Context, ob Object) (*types.Type2Resources, error) {
	if err := s.client.Post(ctx, fmt.Sprintf("/%s", ob.Type), ob.Type2Resources, &Type2); err != nil {
		return nil, fmt.Errorf("failed to add %s %s: %w", ob.Type, ob.Id, err)
	}
	return &Type2Resources, nil
}

func (s *Service) UpdateType2Object(ctx context.Context, ob Object) (*types.Type2Resources, error) {
	if err := s.client.Put(ctx, fmt.Sprintf("/%s/%s", ob.Type, ob.Id), ob.Type2Resources, &Type2); err != nil {
		return nil, fmt.Errorf("failed to update %s %s: %w", ob.Type, ob.Id, err)
	}
	return &Type2Resources, nil
}

func (s *Service) DeleteType2Object(ctx context.Context, ob Object) error {
	if err := s.client.Delete(ctx, fmt.Sprintf("/%s/%s", ob.Type, ob.Id), nil); err != nil {
		return fmt.Errorf("failed to delete %s %s: %w", ob.Type, ob.Id, err)
	}
	return nil
}
