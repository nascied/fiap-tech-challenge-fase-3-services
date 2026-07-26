package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials" //adicionado essa linha para poder autenticar na aws
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// initOTel configura os providers de tracing e métricas do OpenTelemetry, exportando
// via OTLP/gRPC para o OTel Collector (Service dentro do cluster, namespace "monitoring").
// Este é o serviço citado no enunciado para o Distributed Tracing: ele chama flag-service
// e targeting-service via HTTP, então o cliente HTTP também é instrumentado
// (otelhttp.NewTransport, abaixo) para propagar o trace entre os serviços.
func initOTel(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("otel-collector-opentelemetry-collector.monitoring.svc.cluster.local:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("otel-collector-opentelemetry-collector.monitoring.svc.cluster.local:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	return func(ctx context.Context) error {
		return errors.Join(tp.Shutdown(ctx), mp.Shutdown(ctx))
	}, nil
}

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

	// Inicializa o OpenTelemetry (traces → Datadog, métricas → Prometheus)
	shutdown, err := initOTel(ctx, "evaluation-service")
	if err != nil {
		log.Fatalf("erro ao iniciar OpenTelemetry: %v", err)
	}
	defer shutdown(ctx)

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


	// Cliente HTTP (com timeout), com Transport instrumentado para propagar o contexto
	// de trace nas chamadas a flag-service/targeting-service.
	httpClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
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
	if err := http.ListenAndServe(":"+port, otelhttp.NewHandler(mux, "evaluation-service")); err != nil {
		log.Fatal(err)
	}
}