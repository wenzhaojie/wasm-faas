import os
import json
import redis
from PIL import Image
import time
import pickle


# Get Redis URL
def get_url():
    return os.getenv('REDIS_URL', 'redis://127.0.0.1/')


# Store image into Redis
def put_input_img_into_redis(input_img_path, input_obj_key):
    img = Image.open(input_img_path)
    serializable_image = pickle.dumps(img)
    r = redis.Redis.from_url(get_url())
    r.set(input_obj_key, serializable_image)
    print(f"Image saved successfully in Redis with key {input_obj_key}")


# Retrieve image from Redis and save to a path
def get_output_img_from_redis(output_img_path, output_obj_key):
    r = redis.Redis.from_url(get_url())
    data = r.get(output_obj_key)
    image = pickle.loads(data)
    image.save(output_img_path)
    print(f"Image saved successfully at {output_img_path}")


# Apply a simple transformation to the image (like a LUT)
def apply_lut(image):
    processed_image = image.point(lambda p: min(p + 50, 255))
    print(f"Image processed successfully")
    return processed_image


# Handler that uses Redis and measures time
def handler(input_obj_key, output_obj_key):
    stats_dict = {}
    r = redis.Redis.from_url(get_url())

    start = time.time()
    input_obj_data = r.get(input_obj_key)
    stats_dict['get_input_t'] = (time.time() - start) * 1000

    input_obj = pickle.loads(input_obj_data)
    start = time.time()
    output_obj = apply_lut(input_obj)
    stats_dict['compute_t'] = (time.time() - start) * 1000

    start = time.time()
    output_obj_serialized = pickle.dumps(output_obj)
    r.set(output_obj_key, output_obj_serialized)
    stats_dict['set_output_t'] = (time.time() - start) * 1000

    return json.dumps(stats_dict)


def test_pickle():
    img = Image.open('input.jpg')
    img_serialized = pickle.dumps(img)
    img_deserialized = pickle.loads(img_serialized)
    img_deserialized.save('./output.jpg')
    print("Image saved successfully")


def test_put_input_img_into_redis():
    input_img_path = 'input.jpg'
    input_obj_key = 'input_img'
    put_input_img_into_redis(input_img_path, input_obj_key)


def test_apply_lut():
    input_img_path = 'input.jpg'
    input_obj_key = 'input_img'
    put_input_img_into_redis(input_img_path, input_obj_key)
    r = redis.Redis.from_url(get_url())
    input_obj_data = r.get(input_obj_key)
    input_obj = pickle.loads(input_obj_data)
    output_obj = apply_lut(input_obj)
    assert output_obj != None


def test_handler():
    start_time = time.time()
    input_img_path = 'input.jpg'
    input_obj_key = 'input_img'
    output_obj_key = 'output_img'
    put_input_img_into_redis(input_img_path, input_obj_key)
    handler(input_obj_key, output_obj_key)
    get_output_img_from_redis(output_img_path='output.jpg', output_obj_key=output_obj_key)
    end_time = time.time()
    print(f"Time taken: {end_time - start_time}")


if __name__ == '__main__':
    # test_pickle()
    # test_put_input_img_into_redis()
    # test_apply_lut()
    test_handler()