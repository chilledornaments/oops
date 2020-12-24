import requests
import json

from os import environ
from time import sleep

def test_create_secret(tls: bool, url: str, secret: str) -> str:

    j = {"secret": secret}

    r = requests.post(f"{url}/create", verify=tls, json=j)

    if r.status_code != 200:
        raise Exception(f"Status code was {r.status_code}")
    
    secret_url = json.loads(r.text)["url"]

    print("Able to create secret")

    return secret_url
    
def test_get_existing_secret(tls: bool, url: str, secret: str):

    r = requests.get(url, verify=tls)

    if r.status_code != 200:
        raise Exception(f"Status code was {r.status_code}")
    
    if r.text.strip("\n") != secret:
        raise Exception(f"Expected secret to be {secret} but got {r.text}")

    print("Able to view secret")
    

def test_get_already_viewed_secret(tls: bool, url: str):

    r = requests.get(url, verify=tls)

    if r.text.strip("\n") != "Secret not found":
        raise Exception(f"Expected to receive 'Secret not found' message but got {r.text}")

    print("Unable to view already-viewed secret")


def test_expired_secret(tls: bool, url: str):
    
    r = requests.get(url, verify=tls)

    # Response is different when using dynamo
    if environ.get("USING_DYNAMO"):
        if r.text.strip("\n") != "Secret not found":
            raise Exception(f"Expected to receive 'Secret not found' message but got {r.text}")
    else:
        if r.text.strip("\n") != "Secret expired":
            raise Exception(f"Expected to receive 'Secret expired' message but got {r.text}")
    
    print("Unable to view expired secret")

if __name__ == "__main__":
    url = test_create_secret(False, "http://localhost:8080" ,"testing-secret")
    test_get_existing_secret(False, url, "testing-secret")
    test_get_already_viewed_secret(False, url)
    # Create new secret so we can test expiration
    url = test_create_secret(False, "http://localhost:8080" ,"testing-secret")
    if environ.get("USING_DYNAMO"):
        sleep(15)
    else:
        sleep(5)
    test_expired_secret(False, url)

