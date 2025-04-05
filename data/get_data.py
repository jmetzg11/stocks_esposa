import os
import pandas as pd
from dotenv import load_dotenv
import requests
import csv
import os
load_dotenv(dotenv_path='../.env')
import time
from datetime import datetime, timedelta
import sqlite3
from requests.exceptions import ConnectionError, ChunkedEncodingError, ReadTimeout

marketUrl = os.getenv('marketUrl')
url = marketUrl + '/stocks/bars'


class GetData:
    """
    1. get all assests from Alpaca (get_assets)
    2. Get all Market caps from Finhub, watch out for throttling (get_market_cap)
    3. Get 10 years of data of a stock from Alpaca (get_historical_bars)
    4. Put data in sqlite db
    """
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

        # start_here = False
        start_here = True
        for i, row in stocks_df.iterrows():
            symbol = row['symbol']

            # if symbol == 'MIMTF': # last successful operation
            #     start_here = True
            #     continue
            if start_here:
                url = f"https://finnhub.io/api/v1/stock/profile2?symbol={symbol}&token={self.finhub_token}"

                if i % 100 == 0:
                    print(f'On ticker {i}')

                max_retries = 5
                retry_count = 0
                retry = True

                while retry:
                    try:
                        response = requests.get(url, timeout=15)
                        if response.status_code == 429:
                            print("Rate limit hit, pausing for 35 seconds...")
                            time.sleep(35)
                        elif response.status_code == 200:
                            retry = False
                        else:
                            print(f'Issue!!!!!! {response.status_code}')
                            print(response.text)
                            retry_count += 1
                            time.sleep(2)
                    except (ConnectionError, ChunkedEncodingError, ReadTimeout) as e:
                        retry_count += 1
                        wait_time = 5 * retry_count  # Increasing backoff
                        print(f"Connection error: {e}. Retrying in {wait_time} seconds... ({retry_count}/{max_retries})")
                        time.sleep(wait_time)

                if retry_count < max_retries:
                    try:
                        data = response.json()
                        market_cap = data['marketCapitalization']
                        if type(market_cap) == float and market_cap > 100:
                            with open('stocks_with_market_cap.csv', 'a', newline='') as file:
                                writer = csv.writer(file)
                                writer.writerow([symbol, data['marketCapitalization']])

                    except Exception as e:
                        print(f'Problem with stock: {symbol}, error: {e}')
                        with open('errors_fetching.csv', 'a', newline='') as file:
                            writer = csv.writer(file)
                            writer.writerow([symbol])


                    time.sleep(0.5)

                else:
                    print(f"Failed to get data for {symbol} after {max_retries} attempts")
                    with open('errors_fetching.csv', 'a', newline='') as file:
                        writer = csv.writer(file)
                        writer.writerow([symbol])

    def fetch_and_save_historical_data(self, url, symbol):
        all_bars = []
        page_token = None
        max_retries = 3
        retry_count = 0
        while True:
            if page_token:
                full_url = url + f'&page_token={page_token}'
            else:
                full_url = url

            try:
                time.sleep(0.3)
                response = requests.get(full_url, headers=self.headers)

                if response.status_code == 429:
                    sleep_time = 10 + (10 * retry_count)
                    time.sleep(sleep_time)
                    print(f'Rate limit: pausing for {sleep_time} seconds')
                    if retry_count < max_retries:
                        retry_count += 1
                        continue
                    else:
                        with open('errors_getting_historical.csv', 'a', newline='') as file:
                            writer = csv.writer(file)
                            writer.writerow([response.status_code, full_url, response.text])
                        return False

                if response.status_code != 200:
                    with open('errors_getting_historical.csv', 'a', newline='') as file:
                        writer = csv.writer(file)
                        writer.writerow([response.status_code, full_url, response.text])
                    return False

                retry_count = 0
                data = response.json()
                bars = data['bars'][symbol]
                bars = [{'c': b['c'], 't': b['t']} for b in bars]
                all_bars.extend(bars)

                page_token = data.get('next_page_token')
                if not page_token:
                    break
            except Exception as e:
                with open('errors_getting_historical.csv', 'a', newline='') as file:
                    writer = csv.writer(file)
                    writer.writerow(['ERROR', full_url, e])
                return False

        if all_bars:
            df = pd.DataFrame(all_bars)

            filename = f'bars/{symbol}.csv'
            df.to_csv(filename, index=False)
            return True
        else:
            with open('errors_getting_historical.csv', 'a', newline='') as file:
                writer = csv.writer(file)
                writer.writerow(['NO BARS', url, ''])
            return False

    def get_historical_bars(self):
        base_url = self.marketUrl + '/stocks/bars'
        df = pd.read_csv('stocks_with_market_cap.csv')

        today = (datetime.now() - timedelta(days=1)).strftime('%Y-%m-%d')
        ten_years = 365 * 10 + (365//2)
        ten_years_ago = (datetime.now() - timedelta(days=ten_years)).strftime('%Y-%m-%d')
        start = f'&start={ten_years_ago}'
        end = f'&end={today}'

        asof = f'&asof={today}'
        timeframe = '&timeframe=1Day'
        count = 0
        for i, row in df.iterrows():
            if i % 50 == 0:
                print(f'On stock: {i}')

            symbol = row['symbol']
            market_cap = row['marketCapitalization']

            # skip tickers that are numbers
            if not isinstance(symbol, str):
                continue

            # 1 billion
            has_large_market_cap = market_cap > 1000

            # check if it's not a likely OTC/foreign stock
            is_shorter_symbol = len(symbol) < 5
            is_acceptable_5_letter = len(symbol) == 5 and symbol[-1] not in ['F', 'Y']
            is_not_foreign_stock = is_shorter_symbol or is_acceptable_5_letter

            if has_large_market_cap and is_not_foreign_stock:
                url = base_url + f"?symbols={row['symbol']}" + timeframe + start + end + asof

                if self.fetch_and_save_historical_data(url, symbol):
                    print(f'Finished: {symbol}')
                    count += 1
                else:
                    print(f'Error: {symbol}')

    def crate_reference_table(self):
        pass

    def create_historical_table(self):
        conn = sqlite3.connect('stocks.db')
        csv_files = os.listdir('bars')
        print(f'Processing {len(csv_files)} files')
        for i, csv_file in enumerate(csv_files):
            symbol = csv_file.split('.')[0]
            df = pd.read_csv(f'bars/{csv_file}')
            df = df.rename(columns={'c': 'price', 't': 'date'})
            df['date'] = df['date'].str.split('T').str[0]
            df['date'] = pd.to_datetime(df['date'])
            df['symbol'] = symbol

            df.to_sql('historical', conn, if_exists='append', index=False)

            if i % 20 == 0:
                print(f'Finished {i}')

        cursor = conn.cursor()
        cursor.execute('CREATE INDEX idx_symbol ON historical(symbol)')
        cursor.execute('CREATE INDEX idx_date ON historical(date)')
        cursor.execute('CREATE INDEX idx_symbol_date ON historical(symbol, date)')
        conn.commit()
        conn.close()

g = GetData()
g.create_db()
