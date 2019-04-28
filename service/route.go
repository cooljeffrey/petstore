package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/cooljeffrey/petstore/model"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Services struct {
	UserService  UserService
	PetService   PetService
	StoreService StoreService
}

func SetupRoutes(services *Services, logger log.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v2", func(r chi.Router) {
		r.Route("/pet", func(r chi.Router) {
			_ = logger.Log("path", "/pet")
			r.Get("/findByStatus", func(w http.ResponseWriter, r *http.Request) {
				_ = logger.Log("path", "/pet/findByStatus", "method", "get")
				str := r.URL.Query().Get("status")
				params := strings.Split(strings.TrimSpace(str), ",")
				if params != nil && len(params) > 0 {
					pets, err := services.PetService.FindPetsByStatus(r.Context(), params)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
					}
					w.WriteHeader(http.StatusOK)
					_ = encodeResponse(r.Context(), w, pets)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			})

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				_ = logger.Log("path", "/pet", "method", "post")
				var pet *model.Pet
				if e := json.NewDecoder(r.Body).Decode(&pet); e != nil {
					w.WriteHeader(405)
				}
				if err := services.PetService.AddPet(r.Context(), pet); err != nil {
					w.WriteHeader(405)
				}
			})
			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				_ = logger.Log("path", "/pet", "method", "put")
				var pet *model.Pet
				if e := json.NewDecoder(r.Body).Decode(&pet); e != nil {
					w.WriteHeader(405)
				}
				if err := services.PetService.UpdatePet(r.Context(), pet); err != nil {
					w.WriteHeader(http.StatusNotFound)
				}
			})

			r.Route("/{petId}", func(r chi.Router) {
				_ = logger.Log("path", "/pet/{petId}")
				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					_ = logger.Log("path", "/pet/{petId}", "method", "post")
					id, err := strconv.ParseInt(chi.URLParam(r, "petId"), 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
					err = r.ParseForm()
					if err != nil {
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
					err = services.PetService.UpdatePetByID(r.Context(), id, r.PostForm["name"][0], r.PostForm["status"][0])
					if err != nil {
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
					w.WriteHeader(http.StatusOK)
				})
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					id, err := strconv.ParseInt(chi.URLParam(r, "petId"), 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
					}
					pet, err := services.PetService.FindPetByID(r.Context(), id)
					if err != nil || pet == nil {
						w.WriteHeader(http.StatusNotFound)
					}
					err = encodeResponse(r.Context(), w, pet)
					if err != nil {
						_ = level.Error(logger).Log("err", err, "resp", pet)
					}
				})
				r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
					id, err := strconv.ParseInt(chi.URLParam(r, "petId"), 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
					}
					err = services.PetService.DeletePetByID(r.Context(), id)
					if err != nil {
						_ = level.Error(logger).Log("err", err, "petId", id)
					}
					w.WriteHeader(http.StatusNoContent)
				})
				r.Post("/uploadImage", func(w http.ResponseWriter, r *http.Request) {
					idstr := chi.URLParam(r, "petId")
					if idstr == "" {
						w.WriteHeader(http.StatusNotFound)
					}
					id, err := strconv.ParseInt(idstr, 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
					}
					var buf bytes.Buffer
					file, header, err := r.FormFile("file")
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					defer file.Close()
					// Copy the file data to my buffer
					io.Copy(&buf, file)
					err = services.PetService.AddImageUrlForPetByID(r.Context(), id, header.Filename, buf.Bytes())
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
					}
					buf.Reset()
					return
				})
			})
		})

		r.Route("/store", func(r chi.Router) {
			r.Get("/inventory", func(w http.ResponseWriter, r *http.Request) {
				inv, err := services.StoreService.GetInventoriesByStatus(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				err = encodeResponse(r.Context(), w, inv)
				if err != nil {
					_ = level.Error(logger).Log("err", err, "inventory", inv)
				}
			})
			r.Post("/order", func(w http.ResponseWriter, r *http.Request) {
				var order *model.Order
				if e := json.NewDecoder(r.Body).Decode(&order); e != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				o, err := services.StoreService.PlaceOrder(r.Context(), order)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				err = encodeResponse(r.Context(), w, o)
				if err != nil {
					_ = level.Error(logger).Log("err", err, "order", o)
				}
				w.WriteHeader(http.StatusCreated)
			})
			r.Get("/order/{orderId}", func(w http.ResponseWriter, r *http.Request) {
				id, err := strconv.ParseInt(chi.URLParam(r, "orderId"), 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				order, err := services.StoreService.FindOrderByID(r.Context(), id)
				if err != nil || order == nil {
					w.WriteHeader(http.StatusNotFound)
				}
				err = encodeResponse(r.Context(), w, order)
				if err != nil {
					_ = level.Error(logger).Log("err", err, "order", order)
				}
			})

			r.Delete("/order/{orderId}", func(w http.ResponseWriter, r *http.Request) {
				id, err := strconv.ParseInt(chi.URLParam(r, "orderId"), 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				err = services.StoreService.DeleteOrderByID(r.Context(), id)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
				}
				w.WriteHeader(http.StatusNoContent)
			})
		})

		r.Route("/user", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				var user *model.User
				if e := json.NewDecoder(r.Body).Decode(&user); e != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				err := services.UserService.CreateUser(r.Context(), user)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				err = encodeResponse(r.Context(), w, user)
				if err != nil {
					_ = level.Error(logger).Log("err", err, "user", user)
				}
			})
			r.Post("/createWithArray", func(w http.ResponseWriter, r *http.Request) {
				var users []*model.User
				if e := json.NewDecoder(r.Body).Decode(&users); e != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				err := services.UserService.CreateUsersWithArray(r.Context(), users)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				err = encodeResponse(r.Context(), w, users)
				if err != nil {
					_ = level.Error(logger).Log("err", err, "users", users)
				}
			})
			r.Post("/createWithList", func(w http.ResponseWriter, r *http.Request) {
				var users []*model.User
				if e := json.NewDecoder(r.Body).Decode(&users); e != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				err := services.UserService.CreateUsersWithArray(r.Context(), users)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				err = encodeResponse(r.Context(), w, users)
				if err != nil {
					_ = level.Error(logger).Log("err", err, "users", users)
				}
			})

			r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
				username := r.URL.Query().Get("username")
				password := r.URL.Query().Get("password")
				err := services.UserService.Login(r.Context(), username, password)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
				}
				// TODO the following are hardcoded
				w.Header().Set("X-Rate-Limit", "100")
				w.Header().Set("X-Expires-After", "3600")

				w.WriteHeader(http.StatusOK)
			})

			r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
				err := services.UserService.Logout(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
				}
				w.WriteHeader(http.StatusOK)
			})

			r.Route("/{username}", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					username := chi.URLParam(r, "username")
					if username == "" {
						w.WriteHeader(http.StatusBadRequest)
					}
					user, err := services.UserService.GetUserByUsername(r.Context(), username)
					if err != nil || user == nil {
						w.WriteHeader(http.StatusNotFound)
					}
					err = encodeResponse(r.Context(), w, user)
					if err != nil {
						_ = level.Error(logger).Log("err", err, "user", user)
					}
				})
				r.Put("/", func(w http.ResponseWriter, r *http.Request) {
					username := chi.URLParam(r, "username")
					if username == "" {
						w.WriteHeader(http.StatusBadRequest)
					}
					var user *model.User
					if e := json.NewDecoder(r.Body).Decode(&user); e != nil {
						w.WriteHeader(http.StatusBadRequest)
					}
					err := services.UserService.UpdateUserByUsername(r.Context(), username, user)
					if err != nil || user == nil {
						w.WriteHeader(http.StatusNotFound)
					}
					w.WriteHeader(http.StatusOK)
				})
				r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
					username := chi.URLParam(r, "username")
					if username == "" {
						w.WriteHeader(http.StatusBadRequest)
					}
					err := services.UserService.DeleteUserByUsername(r.Context(), username)
					if err != nil {
						_ = level.Error(logger).Log("err", err, "username", username)
					}
					w.WriteHeader(http.StatusNoContent)
				})
			})
		})
	})

	// To serve pet images
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "public")
	fileServer(r, "/images", http.Dir(filesDir))
	return r
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if e, ok := err.(*model.ErrResponse); ok {
		w.WriteHeader(int(e.Code))
	} else {
		w.WriteHeader(500)
	}
	_ = json.NewEncoder(w).Encode(err)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
