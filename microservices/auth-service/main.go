package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	_"github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"

	// ==================== OpenTelemetry (instrumentação - Requisito 3) ====================
	// Descomentar junto com o bloco em initTracer() e no main(), quando o endpoint do
	// OTel Collector e a ferramenta de APM (Datadog/New Relic) estiverem definidos.
	// Dependências a adicionar no go.mod:
	//   go.opentelemetry.io/otel
	//   go.opentelemetry.io/otel/sdk
	//   go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
	//   go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
	// "context"
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// "go.opentelemetry.io/otel/sdk/resource"
	// sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	// =======================================================================================
)

// initTracer configura o TracerProvider do OpenTelemetry, exportando traces via OTLP/gRPC
// para o OTel Collector (Service dentro do cluster, namespace "monitoring").
//
// func initTracer(ctx context.Context, serviceName string) (func(context.Context) error, error) {
// 	exporter, err := otlptracegrpc.New(ctx,
// 		otlptracegrpc.WithEndpoint("otel-collector-opentelemetry-collector.monitoring.svc.cluster.local:4317"),
// 		otlptracegrpc.WithInsecure(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res, _ := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithBatcher(exporter),
// 		sdktrace.WithResource(res),
// 	)
// 	otel.SetTracerProvider(tp)
// 	return tp.Shutdown, nil
// }

// App struct (para injeção de dependência)
type App struct {
	DB         *sql.DB
	MasterKey  string
}

func main() {
	// Carrega o .env para desenvolvimento local. Em produção, isso não fará nada.
	_ = godotenv.Load()

	// Inicializa o OpenTelemetry (descomentar junto com os imports e initTracer acima)
	// shutdown, err := initTracer(context.Background(), "auth-service")
	// if err != nil {
	// 	log.Fatalf("erro ao iniciar OpenTelemetry: %v", err)
	// }
	// defer shutdown(context.Background())

	// --- Configuração ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001" // Porta padrão
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL deve ser definida")
	}

	masterKey := os.Getenv("MASTER_KEY")
	if masterKey == "" {
		log.Fatal("MASTER_KEY deve ser definida")
	}

	// --- Conexão com o Banco ---
	db, err := connectDB(databaseURL)
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	app := &App{
		DB:         db,
		MasterKey:  masterKey,
	}

	// --- Rotas da API ---
	mux := http.NewServeMux()
	mux.HandleFunc("/health", app.healthHandler)

	// Endpoint público para validar uma chave
	mux.HandleFunc("/validate", app.validateKeyHandler)

	// Endpoints de "admin" para criar/gerenciar chaves
	// Eles são protegidos pelo middleware de autenticação
	mux.Handle("/admin/keys", app.masterKeyAuthMiddleware(http.HandlerFunc(app.createKeyHandler)))

	log.Printf("Serviço de Autenticação (Go) rodando na porta %s", port)
	// Com OTel ativo, trocar "mux" por otelhttp.NewHandler(mux, "auth-service")
	// para instrumentar automaticamente todas as rotas HTTP.
	// if err := http.ListenAndServe(":"+port, otelhttp.NewHandler(mux, "auth-service")); err != nil {
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

// connectDB inicializa e testa a conexão com o PostgreSQL
func connectDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Conectado ao PostgreSQL com sucesso!")
	return db, nil
}