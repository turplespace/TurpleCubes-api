package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/turplespace/portos/internal/handlers"
)

// CORS to allow requests from all origins
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetupRoutes() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	fmt.Println(ex)
	path := fmt.Sprintf("%s_web", ex)

	fs := http.FileServer(http.Dir(path))
	http.Handle("/", Cors(fs))

	http.Handle("/api/workspaces", Cors(http.HandlerFunc(handlers.HandleGetWorkspaces)))
	http.Handle("/api/workspace/create", Cors(http.HandlerFunc(handlers.HandleCreateWorkspace)))
	http.Handle("/api/workspace/edit", Cors(http.HandlerFunc(handlers.HandleEditWorkspace)))
	http.Handle("/api/workspace/delete", Cors(http.HandlerFunc(handlers.HandleDeleteWorkspace)))

	http.Handle("/api/workspace/deploy", Cors(http.HandlerFunc(handlers.HandleDeployWorkspace)))
	http.Handle("/api/workspace/redeploy", Cors(http.HandlerFunc(handlers.HandleRedeployWorkspace)))
	http.Handle("/api/workspace/stop", Cors(http.HandlerFunc(handlers.HandleStopWorkspace)))

	http.Handle("/api/cubes", Cors(http.HandlerFunc(handlers.HandleGetCubes)))
	http.Handle("/api/cube/add", Cors(http.HandlerFunc(handlers.HandleAddCubes)))
	http.Handle("/api/cube/edit", Cors(http.HandlerFunc(handlers.HandleEditCube)))
	http.Handle("/api/cube/delete", Cors(http.HandlerFunc(handlers.HandleDeleteCube)))

	http.Handle("/api/cube", Cors(http.HandlerFunc(handlers.HandleGetCubeData)))
	http.Handle("/api/cube/deploy", Cors(http.HandlerFunc(handlers.HandleDeployCube)))
	http.Handle("/api/cube/redeploy", Cors(http.HandlerFunc(handlers.HandleRedeployCube)))
	http.Handle("/api/cube/stop", Cors(http.HandlerFunc(handlers.HandleStopCube)))
	http.Handle("/api/cube/commit", Cors(http.HandlerFunc(handlers.HandleCommitCube)))

	http.Handle("/api/proxy", Cors(http.HandlerFunc(handlers.HandlePostProxy)))

	http.Handle("/api/images", Cors(http.HandlerFunc(handlers.HandleGetImages)))
	http.Handle("/api/logs/stream", Cors(http.HandlerFunc(handlers.HandleLogStream)))

}
