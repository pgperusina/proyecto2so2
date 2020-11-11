from pymongo import MongoClient
import redis
import json
from concurrent import futures
import grpc
import helloworld_pb2
import helloworld_pb2_grpc

redis_id = 1

try:
    client = MongoClient("mongodb://34.67.186.172:27017")
    mongo_db = client.covid19
    print('Connection to monboDB successful')
except Exception as e:
    print('Error connecting to mongoDB - %s', e)

try:
    redis_db = redis.Redis(host="34.66.203.76", port=6379, db=0)
    print('Connection to Redis successful')
except Exception as e:
    print("Error connecting to Redis - %s", e)


class Greeter(helloworld_pb2_grpc.GreeterServicer):
    def SayHello(self, request, context):
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
    try:
        print("***** Insertando caso a Mongo DB *****")
        mongo_db.casos.insert_one(data)
        print("***** Caso insertado en MongoDB: %s", data)
    except Exception as e:
        print("Error posting document to mongo DB")
        print(e)

def insert_redis(data):
    try:
        print("***** Insertando caso a Redis *****")
        global redis_id
        redis_id = redis_id + 1
        redis_db.set(str(redis_id), data)
        print("***** Caso insertado en Redis: %s", data)
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