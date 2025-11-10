# ==================================================================================
#
#       Copyright (c) 2025 Samsung Electronics Co., Ltd. All Rights Reserved.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
# ==================================================================================
import requests
import pytest

BASE_URL = "http://localhost:32007"

def get_model_registration_payload():
    return {
        "modelId": {
            "modelName": "placeholder",
            "modelVersion" : "placeholder"
        },
        "description": "hello world2",
        "modelInformation": {
            "metadata": {
                "author": "someone"
            },
            "inputDataType": "pdcpBytesDl,pdcpBytesUl",
            "outputDataType": "c, d",
            "targetEnvironment": [
                {
                    "platformName": "abc",
                    "environmentType": "env",
                    "dependencyList": "a,b,c"
                }
            ]
        }
    }
    
def submit_model_registration(modelName, modelVersion):
    '''
        Helper function to register model and returns the model_id
    '''
    url = f"{BASE_URL}/ai-ml-model-registration/v1/model-registrations/"
    payload = get_model_registration_payload()
    payload["modelId"]["modelName"] = modelName
    payload["modelId"]["modelVersion"] = modelVersion
    r = requests.post(url=url, json=payload)
    assert r.status_code == 201, f"Job Submission didn't returned 201 but returned {r.status_code}"
    model_id = r.json().get("modelInfo").get("id")
    return model_id


def test_model_registration_and_retrieval():
    '''
        The following test verifies that a model-registration can be successfully submitted and that its model-id is retrievable via the corresponding API.
    '''
    # Submit the Model-registraion job
    modelName = "test-model"
    modelVersion = "10003"
    model_id = submit_model_registration(modelName, modelVersion)
    
    # Check job Status through id
    retrieval_url = f"{BASE_URL}/ai-ml-model-registration/v1/model-registrations/{model_id}"
    r = requests.get(url=retrieval_url)
    response_dict = r.json()
    assert r.status_code == 200, f"Model Retrieval by id didn't returned 200, but returned {r.status_code}"
    assert response_dict.get("id") == model_id, "model_id submitted and retrieved doesn't match"

def test_job_status_retrieval_when_job_not_present():
    '''
        The following test verifies how model-retrieval happens when the model is not registered/present
    '''
    model_id = "invalid"
    expected_message = "record not found"
    # Check job Status through id
    retrieval_url = f"{BASE_URL}/ai-ml-model-registration/v1/model-registrations/{model_id}"
    r = requests.get(url=retrieval_url)
    response_dict = r.json()
    assert r.status_code == 500, f"Model Retrieval by id (when id not present) didn't returned 500, but returned {r.status_code}"
    assert response_dict.get("message") == expected_message, "message(in response) submitted and retrieved doesn't match"


def test_registered_model_deletion():
    '''
        The following test verifies how model-deletion happens when the model is already registered/present.
        The steps include registering the model, retrieving it, deleting it, and then attempting to retrieve it again, which should return an appropriate response.
    '''
    # Submit the Model-registraion job
    modelName = "test-model-deletion"
    modelVersion = "10006"
    model_id = submit_model_registration(modelName, modelVersion)
    
    # Check job Status through id
    retrieval_url = f"{BASE_URL}/ai-ml-model-registration/v1/model-registrations/{model_id}"
    r = requests.get(url=retrieval_url)
    response_dict = r.json()
    assert r.status_code == 200, f"Model Retrieval by id didn't returned 200, but returned {r.status_code}"
    assert response_dict.get("id") == model_id, "model_id submitted and retrieved doesn't match"
    
    # Delete Model
    deletion_url = f"{BASE_URL}/ai-ml-model-registration/v1/model-registrations/{model_id}"
    r = requests.delete(url=deletion_url)
    assert r.status_code == 204, f"Model Deletion by id didn't returned 204, but returned {r.status_code}"
    
    # Check job Status through id (It should not be present)
    retrieval_url = f"{BASE_URL}/ai-ml-model-registration/v1/model-registrations/{model_id}"
    expected_message = "record not found"
    r = requests.get(url=retrieval_url)
    response_dict = r.json()
    response_dict = r.json()
    assert r.status_code == 500, f"Model Retrieval by id (after deletion) didn't returned 500, but returned {r.status_code}"
    assert response_dict.get("message") == expected_message, "message(in response) submitted and retrieved doesn't match"
    
def test_registered_model_deletion_when_model_not_present():
    '''
        The following test verifies how model-deleltion happens when the model is not registered/present
    '''
    model_id = "invalid"
    # Delete Model 
    deletion_url = f"{BASE_URL}/ai-ml-model-registration/v1/model-registrations/{model_id}"
    r = requests.delete(url=deletion_url)
    assert r.status_code == 204, f"Model Deletion by id didn't returned 204, but returned {r.status_code}"
