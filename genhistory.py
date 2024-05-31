import random
from datetime import datetime, timedelta
import sys

def generate_example_data(num_entries):
    data = []
    for _ in range(num_entries):
        timestamp = datetime.now() - timedelta(days=random.randint(0, 365),
                                               hours=random.randint(0, 23),
                                               minutes=random.randint(0, 59),
                                               seconds=random.randint(0, 59))
        ip_address = f"{random.randint(0, 255)}.{random.randint(0, 255)}.{random.randint(0, 255)}.{random.randint(0, 255)}"
        entry = f"{timestamp.strftime('%Y-%m-%d %H:%M:%S')} CEST - {ip_address}"
        data.append(entry)
    return data

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python genhistory.py <num_entries>")
        sys.exit(1)

    num_entries = int(sys.argv[1])
    example_data = generate_example_data(num_entries)

    with open('ip_history.txt', 'w') as file:
        for entry in example_data:
            file.write(entry + '\n')
