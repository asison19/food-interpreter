import pandas as pd
from argparse import ArgumentParser

def display_entries(f):
    print(f)
    df = pd.read_csv(f)

    with pd.option_context('display.max_rows', None, 'display.max_columns', None):
        print(df)


if __name__ == "__main__":
    print("hello")
    parser = ArgumentParser()
    parser.add_argument("-f", "--file", dest='file', help="CSV file from the interpreter of the food entries", metavar="FILE")

    args = parser.parse_args()
    display_entries(args.file)
