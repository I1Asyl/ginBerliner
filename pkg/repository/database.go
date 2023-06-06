package repository

import (
	"fmt"
	"log"

	"github.com/I1Asyl/ginBerliner/models"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
}

type Transaction struct {
	*sqlx.Tx
}

// SetupOrm sets up the database connection
func NewDatabase(dsn string) Database {
	db, err := sqlx.Open("mysql", dsn+"?parseTime=true")
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return Database{db}
}

func (db Database) StartTransaction() Transaction {
	tx := db.MustBegin()
	return Transaction{tx}
}

// func (db Database) (tx Transaction) {
// 	tx.Commit()
// }

func (db Database) GetTeamByTeamName(teamName string) (models.Team, error) {
	var team models.Team
	err := db.Get(&team, "SELECT * FROM team WHERE team_name = ?", teamName)
	return team, err
}

func (db Transaction) GetTeamByTeamName(teamName string) (models.Team, error) {
	var team models.Team
	err := db.Get(&team, "SELECT * FROM team WHERE team_name = ?", teamName)
	return team, err
}

func (db Database) GetUserByUserame(username string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM user WHERE username = ?", username)
	return user, err
}
func (db Transaction) GetUserByUserame(username string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM user WHERE username = ?", username)
	return user, err
}

func (db Database) GetUserTeams(user models.User) ([]models.Team, error) {
	var teams []models.Team
	err := db.Select(&teams, "SELECT * FROM team WHERE team_leader_id = ?", user.Id)
	return teams, err
}
func (db Transaction) GetUserTeams(user models.User) ([]models.Team, error) {
	var teams []models.Team
	err := db.Select(&teams, "SELECT * FROM team WHERE team_leader_id = ?", user.Id)
	return teams, err
}

func (db Database) AddUser(user models.User) error {
	_, err := db.Exec("INSERT INTO user (username, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)", user.Username, user.FirstName, user.LastName, user.Email, user.Password)
	return err
}
func (db Transaction) AddUser(user models.User) error {
	_, err := db.Exec("INSERT INTO user (username, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)", user.Username, user.FirstName, user.LastName, user.Email, user.Password)
	return err
}

func (db Database) AddMembership(membership models.Membership) error {
	_, err := db.Exec("INSERT INTO membership (team_id, user_id, is_editor) VALUES (?, ?, ?)", membership.TeamId, membership.UserId, membership.IsEditor)
	return err
}
func (db Transaction) AddMembership(membership models.Membership) error {
	_, err := db.Exec("INSERT INTO membership (team_id, user_id, is_editor) VALUES (?, ?, ?)", membership.TeamId, membership.UserId, membership.IsEditor)
	return err
}

func (db Database) AddTeam(team models.Team) error {
	_, err := db.Exec("INSERT INTO team (team_name, team_leader_id, team_description) VALUES (?, ?, ?)", team.TeamName, team.TeamLeaderId, team.TeamDescription)
	return err
}
func (db Transaction) AddTeam(team models.Team) error {
	_, err := db.Exec("INSERT INTO team (team_name, team_leader_id, team_description) VALUES (?, ?, ?)", team.TeamName, team.TeamLeaderId, team.TeamDescription)
	return err
}

func (db Database) AddPost(post models.Post) (int, error) {
	var id int
	db.Get(&id, "SELECT `AUTO_INCREMENT` FROM  INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'berliner' AND TABLE_NAME = 'post';")
	_, err := db.Exec("INSERT INTO post (author_type, content, updated_at, created_at) VALUES (?, ?, ?, ?)", post.AuthorType, post.Content, post.UpdatedAt, post.CreatedAt)
	return id, err
}
func (db Transaction) AddPost(post models.Post) error {
	_, err := db.Exec("INSERT INTO post (author_type, content, updated_at, created_at) VALUES (?, ?, ?, ?)", post.AuthorType, post.Content, post.UpdatedAt, post.CreatedAt)
	return err
}

func (db Database) AddUserPost(post models.UserPost) error {
	_, err := db.Exec("INSERT INTO user_post (user_id, post_id) VALUES (?, ?)", post.UserId, post.PostId)
	return err
}
func (db Transaction) AddUserPost(post models.UserPost) error {
	_, err := db.Exec("INSERT INTO user_post (user_id, post_id) VALUES (?, ?)", post.UserId, post.PostId)
	return err
}

func (db Database) AddTeamPost(post models.TeamPost) error {
	_, err := db.Exec("INSERT INTO team_post (team_id, post_id) VALUES (?, ?)", post.TeamId, post.PostId)
	return err
}
func (db Transaction) AddTeamPost(post models.TeamPost) error {
	_, err := db.Exec("INSERT INTO team_post (team_id, post_id) VALUES (?, ?)", post.TeamId, post.PostId)
	return err
}

func (db Database) GetUserPosts(user models.User) ([]models.Post, error) {
	var posts []models.Post
	users := "SELECT following.user_id FROM following WHERE following.follower_id=?"
	userPostsId := fmt.Sprintf("SELECT user_post.post_id FROM user_post WHERE user_post.user_id in (%v)", users)
	err := db.Select(&posts, fmt.Sprintf("SELECT * FROM post WHERE post_id in (%v)", userPostsId), user.Id)
	return posts, err
}

func (db Database) GetTeamPosts(user models.User) ([]models.Post, error) {
	var posts []models.Post
	teams := "SELECT membership.team_id FROM membership WHERE membership.user_id=?"
	err := db.Select(&posts, fmt.Sprintf("SELECT post.* FROM post LEFT JOIN team_post ON team_post.post_id = post.id WHERE team_post.team_id in (%v)", teams), user.Id)
	return posts, err
}

func (db Database) GetFollowing(user models.User) ([]models.User, error) {
	var users []models.User
	err := db.Select(&users, "SELECT * FROM following WHERE follower_id = ?", user.Id)
	return users, err
}

func (db Database) DeleteTeam(team models.Team) error {
	_, err := db.Exec("DELETE FROM team WHERE team_id = ?", team.Id)
	return err
}

func (db Database) AddFollowing(following models.Following) error {
	_, err := db.Exec("INSERT INTO following (follower_id, user_id) VALUES (?, ?)", following.FollowerId, following.UserId)
	return err
}

func (db Database) UpdateTeam(team models.Team) error {
	if team.TeamName != "" {
		_, err := db.Exec("UPDATE team SET team_name = ? WHERE team_id = ?", team.TeamName, team.Id)
		if err != nil {
			return err
		}
	}
	if team.TeamDescription != "" {
		_, err := db.Exec("UPDATE team SET team_description = ? WHERE team_id = ?", team.TeamName, team.Id)
		if err != nil {
			return err
		}
	}
	return nil
}