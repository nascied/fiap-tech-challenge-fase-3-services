import os

from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

def init_tracer(service_name):
    endpoint = os.getenv(
        "OTEL_EXPORTER_OTLP_ENDPOINT",
        "http://localhost:4318/v1/traces"
    )

    provider = TracerProvider(
        resource=Resource.create({
            "service.name": service_name
        })
    )

    exporter = OTLPSpanExporter(endpoint=endpoint)

    provider.add_span_processor(BatchSpanProcessor(exporter))

    trace.set_tracer_provider(provider)
