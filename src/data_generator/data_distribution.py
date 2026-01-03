'''
Docstring for data-generator.data_distribution
Module to generate data distributions for synthetic fiber optic link parameters.
'''

import random

# Sample distribution for fiber lengths
def sample_fiber_length(rng) -> int | float:
  return round(rng.uniform(2.0, 40.0), 2)  # Length in km between 2 and 40 km

# Sample distribution for fiber attenuation db per km
def sample_fiber_att_db_per_km(rng) -> int | float:
  return round(rng.choice([0.2, 0.3, 0.4, 0.5]), 3) # Attenuation between 0.2 and 0.5 dB/km

# Sample distribution for tx power in dbm
def sample_tx_power_dbm(rng) -> int | float:
  return round(rng.choice([2, 4, 6, 8]), 1)  # Tx power in dBm from set values

# Sample distribution for rx sensitivity in dbm
def sample_rx_sensitivity_dbm(rng) -> int | float:
  return round(rng.choice([-27, -28, -29, -20]), 1) # Rx sensitivity in dBm from set values

# Sample distribution for system margin in db
def sample_engineering_margin_db(rng) -> int | float:
  return round(rng.choice([2, 3, 4, 5]), 1)  # System margin in dB from set values

# Sample distribution for splitter loss in db
def sample_splitter_loss_db(rng) -> int | float:
  return round(rng.choice([0.5, 1.0, 1.5, 2.0]), 2)  # Splitter loss in dB from set values

# Sample distribution for number of splices
def sample_num_splices(length_km: float) -> int | float:
  return max(1, int(length_km // 3)) # Approx 1 splice every 3 km

# Sample distribution for splice loss in db
def sample_splice_loss_db(rng) -> int | float:
  return round(rng.uniform(0.05, 0.15), 3)  # Splice loss between 0.05 and 0.15 dB

# Sample distribution for number of connectors
def sample_num_connectors(rng) -> int:
  return rng.choice([2, 4, 6, 8])  # Number of connectors from set values

# Sample distribution for connector loss in db
def sample_connector_loss_db(rng) -> int | float:
  return round(rng.uniform(0.2, 0.5), 3)  # Connector loss between 0.2 and 0.5 dB

# Sample distribution for other loss in db
def sample_other_loss_db(rng) -> int | float:
  return round(rng.uniform(0, 2.5), 2)  # Other losses between 0 and 2.5 dB