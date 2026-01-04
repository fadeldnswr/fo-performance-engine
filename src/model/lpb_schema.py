'''
Docstring for model.lpb_schema
Module defining the schema for Link Power Budget (LPB) analysis.
'''

from pathlib import Path
from typing import Optional
from dataclasses import dataclass

# Define class for LPB analysis
@dataclass
class LPBAnalysis:
  result_csv: Path
  output_dir: Path = Path("./exports/lpb")
  input_links_csv: Optional[Path] = None
  join_key: str = "link_id"
  scenario_col: str = "scenario"
  margin_col: str = "margin_db"
  status_col: str = "lpb_status"
  fiber_length_col: str = "fiber_length_km"
  top_contributors: str = "top_contributor_1"
  worst_n: int = 10
