package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type OktaAuthenticator interface {
	CheckAccess(email string) (*OktaAccessStatus, error)
}

type OktaAuth struct {
	client *okta.Client
	mock   bool
}

type OktaAccessStatus struct {
	VPN        bool
	Production bool
	ConfigTool bool
}

func NewOktaAuth() (*OktaAuth, error) {
	host := os.Getenv("OKTA_HOST")
	token := os.Getenv("OKTA_TOKEN")

	// Enable mock mode if credentials not set
	if host == "" || token == "" {
		fmt.Println("WARNING: Running in Okta mock mode - no credentials provided")
		return &OktaAuth{mock: true}, nil
	}

	_, client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(host),
		okta.WithToken(token),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Okta client: %w", err)
	}

	return &OktaAuth{client: client}, nil
}

func (o *OktaAuth) CheckAccess(email string) (*OktaAccessStatus, error) {
	if o.mock {
		// Mock response - grants all access
		return &OktaAccessStatus{
			VPN:        true,
			Production: true,
			ConfigTool: true,
		}, nil
	}

	user, _, err := o.client.User.GetUser(context.Background(), email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	vpnAccess, err := o.checkGroupMembership(user.Id, "VPN_Access")
	if err != nil {
		return nil, fmt.Errorf("failed to check VPN access: %w", err)
	}

	prodAccess, err := o.checkGroupMembership(user.Id, "Production_Access")
	if err != nil {
		return nil, fmt.Errorf("failed to check Production access: %w", err)
	}

	configToolAccess, err := o.checkGroupMembership(user.Id, "Config_Tool_Access")
	if err != nil {
		return nil, fmt.Errorf("failed to check Config Tool access: %w", err)
	}

	return &OktaAccessStatus{
		VPN:        vpnAccess,
		Production: prodAccess,
		ConfigTool: configToolAccess,
	}, nil
}

func (o *OktaAuth) checkGroupMembership(userId, groupName string) (bool, error) {
	if o.mock {
		return true, nil
	}

	groups, _, err := o.client.User.ListUserGroups(context.Background(), userId)
	if err != nil {
		return false, err
	}

	for _, group := range groups {
		if group.Profile.Name == groupName {
			return true, nil
		}
	}
	return false, nil
}
