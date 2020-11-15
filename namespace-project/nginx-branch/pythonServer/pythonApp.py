from pymongo import MongoClient
import redis
import json
from concurrent import futures
import grpc
import helloworld_pb2
import helloworld_pb2_grpc


from opentelemetry import trace
from opentelemetry.ext import jaeger
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchExportSpanProcessor

trace.set_tracer_provider(TracerProvider())

# create a JaegerSpanExporter
jaeger_exporter = jaeger.JaegerSpanExporter(
    service_name='python-grpc-server',
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
    print('Connection to monboDB successful')
except Exception as e:
    print('Error connecting to mongoDB - %s', e)

try:
    redis_db = redis.Redis(host="34.66.203.76", port=6379, password="sopes1proyecto2", db=0)
    print('Connection to Redis successful')
except Exception as e:
    print("Error connecting to Redis - %s", e)


class Greeter(helloworld_pb2_grpc.GreeterServicer):
    def SayHello(self, request, context):
        tracer = trace.get_tracer(__name__)
        with tracer.start_as_current_span("GRPC-server") as grpc_server:  
            try:
                y = json.loads(request.name)
            except Exception as e:
                print("Error deserializing message - %r", request.name)
                print(e)

            print("request deserialized -- %s", y)
            insert_mongo_db(y)
            insert_redis(request.name)
            return helloworld_pb2.HelloReply(message = '>>> %s' %"Caso insertado en DBs")

def insert_mongo_db(data):
    tracer = trace.get_tracer(__name__)
    with tracer.start_as_current_span("grpc_insert_mongo_db") as grpc_insert_mongo:
        try:
            print("***** Insertando caso a Mongo DB *****")
            mongo_db.casos.insert_one(data)
            print("***** Caso insertado en MongoDB: %s", data)
        except Exception as e:
            print("Error posting document to mongo DB")
            print(e)

def insert_redis(data):
    tracer = trace.get_tracer(__name__)
    with tracer.start_as_current_span("grpc_insert_redis") as grpc_insert_redis:
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

def server():
    print('Starting grpc server...')
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=30))
    helloworld_pb2_grpc.add_GreeterServicer_to_server(Greeter(), server)
    server.add_insecure_port('0.0.0.0:50051')
    server.start()
    print("starting grpc server... ")
    server.wait_for_termination()


# logging.basicConfig(level = logging.INFO, filename = 'python.log', filemode = 'w')
server()