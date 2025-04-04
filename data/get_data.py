import os
import pandas as pd
from dotenv import load_dotenv
import requests
import csv
load_dotenv(dotenv_path='../.env')
import time
from datetime import datetime, timedelta

marketUrl = os.getenv('marketUrl')
url = marketUrl + '/stocks/bars'


class GetData:
    def __init__(self):
        self.marketUrl = os.getenv('marketUrl')
        self.accountUrl = os.getenv('accountUrl')
        self.api_key = os.getenv('key')
        self.api_secret = os.getenv('secret')
        self.finhub_token = os.getenv('finhubAPIToken')
        self.headers = {
            "accept": "application/json",
            'APCA-API-KEY-ID': self.api_key,
            'APCA-API-SECRET-KEY': self.api_secret,
        }

    def get_assets(self):
        """Get all stock symbols"""
        url = self.accountUrl + '/assets'
        response = requests.get(url, headers=self.headers)
        data = response.json()

        symbols = [stock['symbol'] for stock in data]

        with open('stock_symbols.csv', 'w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['symbol'])
            for symbol in symbols:
                writer.writerow([symbol])

    def get_market_cap(self):
        stocks_df = pd.read_csv('stock_symbols.csv')

        # Create the output CSV with headers first
        with open('stocks_with_market_cap.csv', 'w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['symbol', 'marketCapitalization'])

        # Create the errors CSV with header
        with open('errors_fetching.csv', 'w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['symbol'])

        for i, row in stocks_df.iterrows():
            symbol = row['symbol']
            url = f"https://finnhub.io/api/v1/stock/profile2?symbol={symbol}&token={self.finhub_token}"

            if i % 100 == 0:
                print(f'On ticker {i}')

            retry = True
            while retry:
                response = requests.get(url)
                if response.status_code == 429:
                    print("Rate limit hit, pausing for 35 seconds...")
                    time.sleep(35)
                else:
                    retry = False

            try:
                data = response.json()
                with open('stocks_with_market_cap.csv', 'a', newline='') as file:
                    writer = csv.writer(file)
                    writer.writerow([symbol, data['marketCapitalization']])

            except Exception as e:
                print(f'Problem with stock: {symbol}, error: {e}')
                with open('errors_fetching.csv', 'a', newline='') as file:
                    writer = csv.writer(file)
                    writer.writerow([symbol])

            time.sleep(0.5)

    def add_to_db(self, data):
        pass

    def get_historical_bars(self):
        base_url = self.marketUrl + '/bars'
        df = pd.read_csv('stocks_with_market_cap.csv')
        today = datetime.now().strftime('%Y-%m-%d')
        ten_years = 365 * 10
        ten_years_ago = (datetime.now() - timedelta(days=ten_years)).strftime('%Y-%m-%d')

        start = f'&start={today}'
        end = f'&end={ten_years_ago}'

        results = []

        for i, row in df.iterrows():
            symbol = f'?symbols={row['sybmol']}'
            if i < 10:
                MC = row['marketCapitalization']
                MC = MC * 1000000
                if MC > 1000000000:

                    url = base_url + symbol + start + end
                    response = requests.get(url, headers=self.headers)

                    if response.status_code != 200:
                        print(response)

                    else:
                        data = response.json()
                        results.append(data)

        return results





        # symbol,marketCapitalization




g = GetData()
g.get_market_cap()
