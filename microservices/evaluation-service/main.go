package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/aws/aws-sdk-go/aws/credentials" //adicionado essa linha para poder autenticar na aws

	// ==================== OpenTelemetry (instrumentação - Requisito 3) ====================
	// Descomentar junto com o bloco em initTracer() e no main(), quando o endpoint do
	// OTel Collector e a ferramenta de APM (Datadog/New Relic) estiverem definidos.
	// Este é o serviço citado no enunciado para exibir o Distributed Tracing: ele chama
	// flag-service e targeting-service via HTTP, então o cliente HTTP também precisa
	// ser instrumentado (otelhttp.NewTransport) para propagar o trace entre os serviços.
	// Dependências a adicionar no go.mod:
	//   go.opentelemetry.io/otel
	//   go.opentelemetry.io/otel/sdk
	//   go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
	//   go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
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

// Contexto global para o Redis
var ctx = context.Background()

// App struct para injeção de dependência
type App struct {
	RedisClient         *redis.Client
	SqsSvc              *sqs.SQS
	SqsQueueURL         string
	HttpClient          *http.Client
	FlagServiceURL      string
	TargetingServiceURL string
}

func main() {
	_ = godotenv.Load() // Carrega .env para dev local

	// Inicializa o OpenTelemetry (descomentar junto com os imports e initTracer acima)
	// shutdown, err := initTracer(context.Background(), "evaluation-service")
	// if err != nil {
	// 	log.Fatalf("erro ao iniciar OpenTelemetry: %v", err)
	// }
	// defer shutdown(context.Background())

	// --- Configuração ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL deve ser definida (ex: redis://localhost:6379)")
	}

	flagSvcURL := os.Getenv("FLAG_SERVICE_URL")
	if flagSvcURL == "" {
		log.Fatal("FLAG_SERVICE_URL deve ser definida")
	}

	targetingSvcURL := os.Getenv("TARGETING_SERVICE_URL")
	if targetingSvcURL == "" {
		log.Fatal("TARGETING_SERVICE_URL deve ser definida")
	}

	// SQS é opcional no dev local, mas obrigatório em prod
	sqsQueueURL := os.Getenv("AWS_SQS_URL")
	awsRegion := os.Getenv("AWS_REGION")
	if sqsQueueURL == "" {
		log.Println("Atenção: AWS_SQS_URL não definida. Eventos não serão enviados.")
	}
	if awsRegion == "" && sqsQueueURL != "" {
		log.Fatal("AWS_REGION deve ser definida para usar SQS")
	}

	// --- Inicializa Clientes ---
	
	// Cliente Redis
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Não foi possível parsear a URL do Redis: %v", err)
	}
	rdb := redis.NewClient(opt)
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Não foi possível conectar ao Redis: %v", err)
	}
	log.Println("Conectado ao Redis com sucesso!")
// Caso for usar com localstack descomentar esssas linhas
// 	Cliente SQS (AWS SDK)
// 	var sqsSvc *sqs.SQS
// 	if sqsQueueURL != "" {
// 		sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion), Endpoint: aws.String(sqsQueueURL)})
// 		if err != nil {
// 			log.Fatalf("Não foi possível criar sessão AWS: %v", err)
// 		}
// 		sqsSvc = sqs.New(sess)
// 		log.Println("Cliente SQS inicializado com sucesso.")
// 	}
// 	var sqsSvc *sqs.SQS
//     if sqsQueueURL != "" {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(awsRegion),
// 	})
// 	if err != nil {
// 		log.Fatalf("Não foi possível criar sessão AWS: %v", err)
// 	}

// 	sqsSvc = sqs.New(sess)
// 	log.Println("Cliente SQS inicializado com sucesso.")
    
//    }

//É possível usar essa linha de baixo basta forncer para localstack qualquer string de autenticação
	   var sqsSvc *sqs.SQS
	   if sqsQueueURL != "" {
		sqsEndpoint := os.Getenv("AWS_ENDPOINT_URL")
		
		awsConfig := &aws.Config{
			Region: aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				os.Getenv("AWS_SESSION_TOKEN"),
			),
		}
		
		if sqsEndpoint != "" {
			log.Printf("Usando SQS local: %s", sqsEndpoint)
			awsConfig.Endpoint = aws.String(sqsEndpoint)
		} else {
			log.Println("Usando SQS da AWS")
		}
		
		sess, err := session.NewSession(awsConfig)
		if err != nil {
			log.Fatalf("Não foi possível criar sessão AWS: %v", err)
		}
		sqsSvc = sqs.New(sess)
		log.Println("Cliente SQS inicializado com sucesso.")
	}


	// Cliente HTTP (com timeout)
	// Com OTel ativo, adicionar Transport: otelhttp.NewTransport(http.DefaultTransport)
	// para propagar o contexto de trace nas chamadas a flag-service/targeting-service.
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		// Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Cria a instância da App
	app := &App{
		RedisClient:         rdb,
		SqsSvc:              sqsSvc,
		SqsQueueURL:         sqsQueueURL,
		HttpClient:          httpClient,
		FlagServiceURL:      flagSvcURL,
		TargetingServiceURL: targetingSvcURL,
	}

	// --- Rotas ---
	mux := http.NewServeMux()
	mux.HandleFunc("/health", app.healthHandler)
	mux.HandleFunc("/evaluate", app.evaluationHandler)

	log.Printf("Serviço de Avaliação (Go) rodando na porta %s", port)
	// Com OTel ativo, trocar "mux" por otelhttp.NewHandler(mux, "evaluation-service")
	// para instrumentar automaticamente todas as rotas HTTP.
	// if err := http.ListenAndServe(":"+port, otelhttp.NewHandler(mux, "evaluation-service")); err != nil {
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}