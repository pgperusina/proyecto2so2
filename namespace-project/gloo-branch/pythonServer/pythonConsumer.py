import pika
from pymongo import MongoClient
import redis
import json

from opentelemetry import trace
from opentelemetry.ext import jaeger
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchExportSpanProcessor

trace.set_tracer_provider(TracerProvider())

# create a JaegerSpanExporter
jaeger_exporter = jaeger.JaegerSpanExporter(
    service_name='python-rabbitmq-consumer',
    # configure agent
    agent_host_name='jaeger-agent.observability.svc.cluster.local',
    agent_port=6831,
    # optional: configure also collector
    # collector_host_name='localhost',
    # collector_port=14268,
    # collector_endpoint='/api/traces?format=jaeger.thrift',
    # username=xxxx, # optional
    # password=xxxx, # optional
)

# Create a BatchExportSpanProcessor and add the exporter to it
span_processor = BatchExportSpanProcessor(jaeger_exporter)

# add to the tracer
trace.get_tracer_provider().add_span_processor(span_processor)

redis_id = 1
KEY_INDEX = 'index'

try:
    client = MongoClient("mongodb://sopes1:sopes1proyecto2@34.67.186.172:27017")
    mongo_db = client.covid19
except Exception as e:
    print('Error connecting to mongoDB - %s', e)

try:
    redis_db = redis.Redis(host="34.66.203.76", port=6379, password="sopes1proyecto2", db=0)
except Exception as e:
    print("Error connecting to Redis - %s", e)


def consume_queue():
    print("tracing consume_rabbitmq_queue")
    credentials = pika.PlainCredentials('guest', 'guest')
    parameters = pika.ConnectionParameters("34.72.226.148", 5672)
    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()
    channel.queue_declare(queue='covid19')
    channel.basic_consume(queue='covid19', on_message_callback=persistInDB, auto_ack=True)
    print('Waiting for messages...')
    channel.start_consuming()

def persistInDB(ch, method, properties, body):
    tracer = trace.get_tracer(__name__)
    with tracer.start_as_current_span("db-persist") as persist_in_db_span:
        print("tracing persistInDB")

        print("Covid19 case received: " % body)

        try:
            data = json.loads(body)
        except Exception as e:
            print("Error deserializing message - %r", body)
            print(e)
        insert_mongo_db(data)
        insert_redis(body)

def insert_mongo_db(data):
    tracer = trace.get_tracer(__name__)
    with tracer.start_as_current_span("insert_mongo_db") as insert_mongo_db_span:
        print("tracing insert_mongo_db")
        try:
            print("***** Insertando caso a Mongo DB *****")
            mongo_db.casos.insert_one(data)
            print("***** Caso insertado en MongoDB: ", data)
        except Exception as e:
            print("Error posting document to mongo DB")
            print(e)


def insert_redis(data):
    tracer = trace.get_tracer(__name__)
    with tracer.start_as_current_span("insert_redis") as insert_redis_span:
        print("tracing insert_redis")
        try:
            print("***** Insertando caso a Redis *****")
            redis_db.incr(KEY_INDEX, 1)
            index = redis_db.get(KEY_INDEX).decode('utf-8')
            int_index = int(index) 
            redis_db.set(str(int_index), data)
            print("***** Caso insertado en Redis: ", data)
        except Exception as e:
            print("Error posting document to Redis")
            print(e)

consume_queue()