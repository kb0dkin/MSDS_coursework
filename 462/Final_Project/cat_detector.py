#!/usr/bin/env python3
#
#


from subprocess import run
from PIL import Image as PIL_Img
import tflite_runtime.interpreter as tflite
from datetime import datetime as dt
import numpy as np

# this will capture an image using the libcamera library every 1 second (minimum I could get it to do)
# I then read that jpg in, use PIL to convert it into a correctly sized numpy array, then predict using
# the loaded model

# for info on all of this, look through the documentation on the tensorflow website:
# 
# Loading the model: tensorflow.org/lite/guide/python
# Running in python: tensorflow.org/lite/guide/inference#load_and_run_a_model_in_python



linesplit = '---------------------------------------'
interpreter = tflite.Interpreter(model_path='/home/pi/cat_tagger_model.tflite')
my_signature = interpreter.get_signature_runner()

while True:
    image_name = f"cat_jpgs/detect{dt.now().strftime('%Y%m%d_%H:%M:%S')}.jpg" # can change if we want
    run(['libcamera-jpeg','-n','-v','0','--immediate','--timelapse=1','-o',image_name])
    img = PIL_Img.open(image_name).convert('RGB').resize((224,224)) # pull in the image -- resize properly etc
    ary = np.expand_dims(np.asarray(img), axis=0).astype(np.float32) # # convert to a numpy array
    
    print(linesplit)
    output = my_signature(input_1=ary) # make the prediction
    probs = output['dense_1'][0,0]
    print(probs)

    del(img)




