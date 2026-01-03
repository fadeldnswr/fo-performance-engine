'''
Docstring for src.data-generator.generate_links
Module to generate synthetic link data for testing and analysis.
'''
import random
import csv
import argparse

from . import data_distribution as dst
from pathlib import Path

# Define class for generate synthetic link data
class LinkGenerator:
  # Define CSV header columns
  _HEADERS = [
    "link_id", "scenario", "tx_power_dbm", "rx_sensitivity_dbm", "engineering_margin_db",
    "fiber_length_km", "fiber_att_db_per_km", "n_splice", "splice_loss_db",
    "n_connector", "connector_loss_db", "splitter_loss_db", "other_loss_db",
  ]
  def __init__(self, n_links: int, output_csv: Path, seed:int | None = 42, scenario: str = "base"):
    self.n_links = n_links
    self.output_csv = output_csv
    self.seed = seed
    self.scenario = scenario
    self.rng = random.Random(seed)
  def generate_links(self) -> None:
    '''
    Docstring for generate_links
    :param self: Description
    '''
    try:
      # Create directory if it doesn't exist
      if not self.output_csv.parent.exists():
        self.output_csv.parent.mkdir(parents=True, exist_ok=True)
      
      # Write header cols to CSV
      with open(self.output_csv, mode='w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(self._HEADERS)
        
        # Iterate to generate each link
        for i in range(1, self.n_links + 1):
          length_km = dst.sample_fiber_length(self.rng)
          row = [
            f"link_{i:05d}",
            self.scenario,
              dst.sample_tx_power_dbm(self.rng),
              dst.sample_rx_sensitivity_dbm(self.rng),
              dst.sample_engineering_margin_db(self.rng),
              length_km,
              dst.sample_fiber_att_db_per_km(self.rng),
              dst.sample_num_splices(length_km),
              dst.sample_splice_loss_db(self.rng),
              dst.sample_num_connectors(self.rng),
              dst.sample_connector_loss_db(self.rng),
              dst.sample_splitter_loss_db(self.rng),
              dst.sample_other_loss_db(self.rng),
          ]
          # Write row to CSV
          writer.writerow(row)
      print(f"DONE — {self.n_links} links written to {self.output_csv}")
    except Exception as e:
      print(f"An error has occurred due to {e}")
  

# Define function to parse command line arguments
def parse_synthetic_data() -> None:
  '''
  Docstring for parse_args
  '''
  # Argument parser for command line execution
  parser = argparse.ArgumentParser(description="Generate synthetic fiber optic link data.")
  parser.add_argument('--n', type=int, default=1000, help='Number of links to generate')
  parser.add_argument('--output', type=Path, default=Path("./examples/links_generated.csv"), help="Output CSV file path")
  parser.add_argument('--seed', type=int, default=42, help="Random seed for reproducibility")
  parser.add_argument('--scenario', type=str, default="base", help="Scenario label for the generated links")
  args = parser.parse_args()
  
  # Call link generation
  generator = LinkGenerator(
    n_links=args.n,
    output_csv=args.output,
    seed=args.seed,
    scenario=args.scenario
  )
  generator.generate_links()

if __name__ == "__main__":
  parse_synthetic_data()