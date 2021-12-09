package github

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

type LastRepo struct {
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c LastRepo) String() string {
	msg := fmt.Sprintf("\nThe last repo updated \n\n **%v** \n \t %v \n Created At: %s \n Updated At: %s", c.Name, c.Description, c.CreatedAt.Format("2006-01-02"), c.UpdatedAt.Format("2006-01-02"))
	if len(c.Name) == 0 {
		msg = fmt.Sprintf("Could not find a repository")
	}

	return msg
}

func (g *GithubService) GetLastRepoByUsername(ctx context.Context, username string) (*LastRepo, error) {

	options := GetContributionsByUsernameOptions{
		Username: username,
	}

	lastRepo := &LastRepo{
		Name:        "",
		Description: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	var lastRepoQuery struct {
		User struct {
			Repositories struct {
				Nodes []struct {
					Name        string
					Description string
					CreatedAt   time.Time
					UpdatedAt   time.Time
				}
			} `graphql:"repositories(last: 1, ownerAffiliations: OWNER, isFork: false, orderBy: {field: UPDATED_AT, direction: ASC} )"`
		} `graphql:"user(login: $username)"`
	}

	err := g.githubClient.Query(ctx, &lastRepoQuery, map[string]interface{}{
		"username": githubv4.String(options.Username),
	})
	if err != nil {
		return nil, errors.Wrap(err, "github client query")
	}
	lastRepo.Name = lastRepoQuery.User.Repositories.Nodes[0].Name
	lastRepo.Description = lastRepoQuery.User.Repositories.Nodes[0].Description
	lastRepo.CreatedAt = lastRepoQuery.User.Repositories.Nodes[0].CreatedAt
	lastRepo.UpdatedAt = lastRepoQuery.User.Repositories.Nodes[0].UpdatedAt
	return lastRepo, nil

}
