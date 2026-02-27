package main

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jesusthecreator017/fswithgo/cmd/api/helpers"
	"github.com/jesusthecreator017/fswithgo/internal/auth"
	"github.com/jesusthecreator017/fswithgo/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) registerUserHandler(w http.ResponseWriter, req *http.Request) {
	// Create the input struct
	var input struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// Read the json
	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create the errors map
	errs := make(map[string]string)

	// Use an email regex to make sure the email is valid
	input.Email = strings.TrimSpace(input.Email)
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	re := regexp.MustCompile(emailRegex)

	// Validate the email
	if input.Email == "" {
		errs["email"] = "must enter an email"
	} else if !re.MatchString(input.Email) {
		errs["email"] = "not a valid email"
	}

	// Validate the name
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		errs["name"] = "must enter a name"
	} else if len(input.Name) > 255 {
		errs["name"] = "must be shorter than 255 character"
	}

	// Validate the password
	input.Password = strings.TrimSpace(input.Password)
	if input.Password == "" {
		errs["password"] = "must enter a password"
	} else {
		var pwErrs []string
		if len(input.Password) < 8 {
			pwErrs = append(pwErrs, "at least 8 characters")
		}
		if len(input.Password) > 20 {
			pwErrs = append(pwErrs, "20 characters or fewer")
		}
		if !regexp.MustCompile(`[A-Z]`).MatchString(input.Password) {
			pwErrs = append(pwErrs, "one uppercase letter")
		}
		if !regexp.MustCompile(`[a-z]`).MatchString(input.Password) {
			pwErrs = append(pwErrs, "one lowercase letter")
		}
		if !regexp.MustCompile(`[0-9]`).MatchString(input.Password) {
			pwErrs = append(pwErrs, "one number")
		}
		if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(input.Password) {
			pwErrs = append(pwErrs, "one special character (!@#$%^&*)")
		}
		if len(pwErrs) > 0 {
			errs["password"] = "must contain: " + strings.Join(pwErrs, ", ")
		}
	}

	// Check the errors
	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create the User Store object
	user := &store.User{
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: string(hashedPassword),
	}

	if err := app.store.Users.Create(req.Context(), user); err != nil {
		// Check if its a duplicate email
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			helpers.ErrorJson(w, http.StatusConflict, "email already taken")
			return
		}

		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// Generate the JWT Token
	token, err := auth.GenerateJWT(user.ID, user.Permissions, app.config.jwtSecret, time.Hour*24*7)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create jwt token")
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelope{"user": user, "token": token})
}

func (app *application) loginUserHandler(w http.ResponseWriter, req *http.Request) {
	// Parase the JSON body
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "error reading JSON body")
		return
	}

	// Create errors map
	errs := make(map[string]string)

	input.Email = strings.TrimSpace(input.Email)
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	re := regexp.MustCompile(emailRegex)
	if input.Email == "" {
		errs["email"] = "must enter an email"
	} else if !re.MatchString(input.Email) {
		errs["email"] = "not a valid email"
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	user, err := app.store.Users.GetByEmail(req.Context(), input.Email)
	if err != nil {
		helpers.ErrorJson(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Compare password with bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		helpers.ErrorJson(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Generate the JWT Token
	token, err := auth.GenerateJWT(user.ID, user.Permissions, app.config.jwtSecret, time.Hour*24*7)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "error generating JWT token")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"user": user, "token": token})
}
