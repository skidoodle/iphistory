import random
import json
from datetime import datetime
import os
import sys

IP_HISTORY_FILE = 'history.json'

def generate_ipv4():
    return '.'.join(str(random.randint(0, 255)) for _ in range(4))

def generate_ipv6():
    return ':'.join('{:x}'.format(random.randint(0, 65535)) for _ in range(8))

def generate_timestamp():
    return datetime.now().strftime("%Y-%m-%d %H:%M:%S %Z")

def generate_ip_record():
    return {
        "timestamp": generate_timestamp(),
        "ipv4": generate_ipv4(),
        "ipv6": generate_ipv6()
    }

def load_ip_history():
    if os.path.exists(IP_HISTORY_FILE):
        with open(IP_HISTORY_FILE, 'r') as file:
            try:
                return json.load(file)
            except json.JSONDecodeError:
                return []
    return []

def save_ip_history(ip_history):
    with open(IP_HISTORY_FILE, 'w') as file:
        json.dump(ip_history, file, indent=2)

def generate_and_save_ip_records(num_records=1):
    ip_history = load_ip_history()
    new_records = [generate_ip_record() for _ in range(num_records)]
    ip_history = new_records + ip_history  # Append new records at the beginning
    save_ip_history(ip_history)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python genhistory.py <number_of_records>")
        sys.exit(1)

    try:
        num_records = int(sys.argv[1])
        if num_records <= 0:
            raise ValueError("Number of records must be a positive integer.")

        generate_and_save_ip_records(num_records)
        print(f"Successfully generated {num_records} IP records.")

    except ValueError as e:
        print(f"Error: {e}")
        print("Please provide a valid positive integer for the number of records.")
        sys.exit(1)
