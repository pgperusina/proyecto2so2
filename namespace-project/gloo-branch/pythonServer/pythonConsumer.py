import pika
from pymongo import MongoClient
import redis
import json

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
    credentials = pika.PlainCredentials('guest', 'guest')
    parameters = pika.ConnectionParameters("34.72.226.148", 5672)
    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()
    channel.queue_declare(queue='covid19')
    channel.basic_consume(queue='covid19', on_message_callback=persistInDB, auto_ack=True)
    print('Waiting for messages...')
    channel.start_consuming()

def persistInDB(ch, method, properties, body):
    print("Covid19 case received: %r" % body)

    try:
        data = json.loads(body)
    except Exception as e:
        print("Error deserializing message - %r", body)
        print(e)
    insert_mongo_db(data)
    insert_redis(body)

def insert_mongo_db(data):
    try:
        print("***** Insertando caso a Mongo DB *****")
        mongo_db.casos.insert_one(data)
        print("***** Caso insertado en MongoDB: ", data)
    except Exception as e:
        print("Error posting document to mongo DB")
        print(e)

def insert_redis(data):
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