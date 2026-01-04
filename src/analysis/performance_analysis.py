'''
Docstring for analysis.lpb_analysis
Module to analyze link power budgets from synthetic link data.
'''
from __future__ import annotations

import csv
import argparse
import sys
import matplotlib.pyplot as plt
import pandas as pd

from typing import Tuple, Optional
from pathlib import Path

from src.exception.exception import CustomException
from src.model.lpb_schema import LPBAnalysis
from src.logging.logging import logging

# Define function to load results CSV
def load_results(config:LPBAnalysis) -> pd.DataFrame:
  '''
  Load results from CSV file specified in config.
  
  :param config: Configuration for LPB analysis
  :type config: LPBAnalysis
  :return: DataFrame of results
  :rtype: pd.DataFrame
  '''
  try:
    # Check if results CSV exists
    if not config.result_csv.exists():
      logging.error(f"Result CSV file not found: {config.result_csv}")
      raise FileNotFoundError(f"Result CSV file not found: {config.result_csv}")
    
    # Load CSV into DataFrame
    df = pd.read_csv(config.result_csv)
    
    # Perform basic sanity checks
    required_cols = {config.join_key, config.scenario_col, config.margin_col, config.status_col}
    missing_cols = [col for col in required_cols if col not in df.columns]
    if missing_cols:
      logging.error(f"Missing required columns in results CSV: {missing_cols}")
      raise ValueError(f"Missing required columns in results CSV: {missing_cols}")
    
    # Normalize column data types
    df[config.margin_col] = pd.to_numeric(df[config.margin_col], errors='coerce')
    df[config.status_col] = df[config.status_col].astype(str).str.upper().str.strip()
    df[config.scenario_col] = df[config.scenario_col].astype(str)
    return df
  except Exception as e:
    raise CustomException(e, sys)

# Define function to merge inputs
def merge_inputs(config:LPBAnalysis, df:pd.DataFrame) -> pd.DataFrame:
  '''
  Merge input link parameters into results DataFrame.
  
  :param config: Configuration for LPB analysis
  :type config: LPBAnalysis
  :param df: Results DataFrame
  :type df: pd.DataFrame
  :return: Merged DataFrame with input link parameters
  :rtype: pd.DataFrame
  '''
  try:
    # Check if input links CSV is provided
    if config.input_links_csv is None:
      logging.info("No input links CSV provided; skipping merge.")
      return df
    # Check if input links CSV exists
    if not config.input_links_csv.exists():
      logging.error(f"Input links CSV file not found: {config.input_links_csv}")
      raise FileNotFoundError(f"Input links CSV file not found: {config.input_links_csv}")
    
    # Load input links CSV
    input_df = pd.read_csv(config.input_links_csv)
    if config.join_key not in input_df.columns:
      logging.error(f"Join key '{config.join_key}' not found in input links CSV.")
      raise ValueError(f"Join key '{config.join_key}' not found in input links CSV.")
    
    # Merge dataframes on join key
    cols_to_merge = [config.join_key]
    for col in [config.fiber_length_col, "splitter_loss_db", "fiber_att_db_per_km", "engineering_margin_db"]:
      if col in input_df.columns:
        cols_to_merge.append(col)
    
    # Copy the cols to merge to avoid SettingWithCopyWarning
    inp = input_df[cols_to_merge].copy()
    merged = df.merge(inp, on=config.join_key, how='left', suffixes=('', '_in'))
    return merged
  except Exception as e:
    raise CustomException(e, sys)

# Define function summarize data
def summarize(df:pd.DataFrame, config:LPBAnalysis) -> pd.DataFrame:
  '''
  Summarize LPB analysis results and export it into CSV.
  
  :param df: Input dataframe
  :type df: pd.DataFrame
  :param config: Configuration for LPB analysis
  :type config: LPBAnalysis
  :return: Summary DataFrame
  :rtype: pd.DataFrame
  '''
  try:
    # Group by scenario column
    group_df = df.groupby(config.scenario_col, dropna=False)
    
    # Perform fail and pass definition
    pass_rate = group_df.apply(lambda x: (x[config.status_col] == "PASS").mean())
    fail_rate = 1 - pass_rate
    
    # Create summary dataframe
    summary_df = pd.DataFrame({
      "n_links": group_df.size(),
      "pass_rate": pass_rate,
      "fail_rate": fail_rate,
      "margin_mean_db": group_df[config.margin_col].mean(),
      "margin_median_db": group_df[config.margin_col].median(),
      "margin_p05_db": group_df[config.margin_col].quantile(0.05),
      "margin_p95_db": group_df[config.margin_col].quantile(0.95),
      "margin_min_db": group_df[config.margin_col].min(),
      "margin_max_db": group_df[config.margin_col].max(),
    }).reset_index()
    
    # Sort dataframe by fail rate descending, then margin mean ascending
    summary = summary_df.sort_values(by=["fail_rate", "margin_mean_db"], ascending=[False, True])
    return summary
  except Exception as e:
    raise CustomException(e, sys)

# Define function to save summary into tables
def save_summary(df:pd.DataFrame, config:LPBAnalysis) -> Tuple[Path, Path]:
  '''
  Save summary DataFrame into CSV files.
  
  :param df: Input dataframe
  :type df: pd.DataFrame
  :param config: LPB analysis configuration
  :type config: LPBAnalysis
  :return: Paths to summary CSV files
  :rtype: Tuple[Path, Path]
  '''
  try:
    # Create output directory if it doesn't exist
    if not config.output_dir.exists():
      config.output_dir.mkdir(parents=True, exist_ok=True)
    
    # Define summary dataframe
    summary_df = summarize(df, config=config)
    summary_path = config.output_dir / "summary_by_scenario.csv"
    summary_df.to_csv(summary_path, index=False)
    
    # Define worst scenario dataframe
    worst = df[[config.join_key, config.scenario_col, config.margin_col, config.status_col]].copy()
    worst = worst.sort_values(by=config.margin_col, ascending=True).head(config.worst_n)
    worst_path = config.output_dir / f"worst_{config.worst_n}.csv"
    worst.to_csv(worst_path, index=False)
    
    return summary_path, worst_path
  except Exception as e:
    raise CustomException(e, sys)

# Define function to plot margin distribution
def plot_margin_hist(df:pd.DataFrame, config:LPBAnalysis) -> Path:
  '''
  Plot histogram of margin distribution.
  
  :param df: Input dataframe
  :type df: pd.DataFrame
  :param config: LPB analysis configuration
  :type config: LPBAnalysis
  :return: Path to saved plot
  :rtype: Path
  '''
  try:
    # Check if output directory exists
    if not config.output_dir.exists():
      config.output_dir.mkdir(parents=True, exist_ok=True)
    
    # Plot histogram
    fig, ax = plt.subplots()
    series = df[config.margin_col].dropna()
    ax.hist(series, bins=30, color='blue', alpha=0.7)
    ax.axvline(0.0, color='red', linestyle='dashed', linewidth=1)
    
    # Define labels and title
    ax.set_xlabel('Margin (dB)')
    ax.set_ylabel('Number of Links')
    ax.set_title('Link Power Budget Margin Distribution')
    
    # Save plot to file
    output_path = config.output_dir / "margin_distribution.png"
    fig.savefig(output_path, dpi=200, bbox_inches="tight")
    plt.close(fig)
    return output_path
  except Exception as e:
    raise CustomException(e, sys)

# Define function to plot pass or fail chart
def plot_pass_or_fail(df:pd.DataFrame, config:LPBAnalysis) -> Path:
  '''
  Plot pass or fail bar chart for different link scenario.
  
  :param df: Input dataframe
  :type df: pd.DataFrame
  :param config: LPB analysis configuration
  :type config: LPBAnalysis
  :return: Path to saved plot
  :rtype: Path
  '''
  try:
    # Check if output directory exists
    if not config.output_dir.exists():
      config.output_dir.mkdir(parents=True, exist_ok=True)
    
    # Plot pass or fail bar chart
    fig, ax = plt.subplots()
    counts = df[config.status_col].value_counts() # Count pass and fail
    order = [x for x in ["PASS", "FAIL"] if x in counts.index] + [x for x in counts.index if x not in ["PASS", "FAIL"]] # Define order of bars
    counts = counts.reindex(order) # Reindex counts to match order
    
    # Create bar chart
    ax.bar(counts.index.astype(str), counts.values)
    ax.set_title('Link Power Budget Pass or Fail Summary')
    ax.set_xlabel('Status')
    ax.set_ylabel('Number of Links')
    
    # Save plot to file
    output_path = config.output_dir / "pass_or_fail_summary.png"
    fig.savefig(output_path, dpi=200, bbox_inches="tight")
    plt.close(fig)
    return output_path
  except Exception as e:
    raise CustomException(e, sys)

# Define function to plot margin vs length
def plot_margin_vs_length(df:pd.DataFrame, config:LPBAnalysis) -> Optional[Path]:
  '''
  Plot margin vs fiber length scatter plot.
  
  :param df: Input dataframe
  :type df: pd.DataFrame
  :param config: LPB analysis configuration
  :type config: LPBAnalysis
  :return: Path to saved plot
  :rtype: Path | None
  '''
  try:
    # Check if fiber length column exists
    if config.fiber_length_col not in df.columns:
      logging.warning(f"Fiber length column '{config.fiber_length_col}' not found; skipping margin vs length plot.")
      return None
    
    # Check if output directory exists
    if not config.output_dir.exists():
      config.output_dir.mkdir(parents=True, exist_ok=True)
    
    # Plot margin vs fiber length scatter plot
    fig, ax = plt.subplots()
    x = pd.to_numeric(df[config.fiber_length_col], errors='coerce')
    y = df[config.margin_col]
    
    # Simple coloring for indicatin pass or fail
    is_pass = df[config.status_col] == "PASS"
    ax.scatter(x[is_pass], y[is_pass], s=10, label="PASS")
    ax.scatter(x[~is_pass], y[~is_pass], s=10, label="FAIL")
    
    # Define labels and title
    ax.axhline(0.0, color='red', linestyle='dashed', linewidth=1)
    ax.set_title('Margin vs Fiber Length')
    ax.set_xlabel('Fiber Length (km)')
    ax.set_ylabel('Margin (dB)')
    ax.legend()
    
    # Save plot to file
    output_path = config.output_dir / "margin_vs_fiber_length.png"
    fig.savefig(output_path, dpi=200, bbox_inches="tight")
    plt.close(fig)
    return output_path
  except Exception as e:
    raise CustomException(e, sys)

# Define function to plot top contributors
def plot_top_contributors(df:pd.DataFrame, config:LPBAnalysis) -> Optional[Path]:
  '''
  Plot top contributors to margin failures.
  
  :param df: Input dataframe
  :type df: pd.DataFrame
  :param config: LPB analysis configuration
  :type config: LPBAnalysis
  :return: Path to saved plot
  :rtype: Path | None
  '''
  try:
    # Check for top contributor column
    if config.top_contributors not in df.columns:
      logging.warning(f"Top contributor column '{config.top_contributors}' not found; skipping top contributors plot.")
      return None
    
    # Check if output directory exists
    if not config.output_dir.exists():
      config.output_dir.mkdir(parents=True, exist_ok=True)
    
    # Plot top contributors bar chart
    fig, ax = plt.subplots()
    counts = df[config.top_contributors].astype(str).value_counts().head(10) # Top 10 contributors
    
    ax.barh(counts.index[::-1], counts.values[::-1])  # Horizontal bar chart
    ax.set_title('Top 10 Contributors to Margin Failures')
    ax.set_xlabel('Number of Links')
    ax.set_ylabel(f'Contributor: {config.top_contributors}')
    
    # Save plot to file
    output_path = config.output_dir / "top_contributors.png"
    fig.savefig(output_path, dpi=200, bbox_inches="tight")
    plt.close(fig)
    return output_path
  except Exception as e:
    raise CustomException(e, sys)

# Define main function to run the analysis
def run_analysis(config: LPBAnalysis) -> None:
  '''
  Docstring for run_analysis
  
  :param config: Description
  :type config: LPBAnalysis
  '''
  try:
    # Define dataframe
    df = load_results(config=config)
    df = merge_inputs(config=config, df=df)
    
    # Save summary tables
    summary_path, worst_path = save_summary(df=df, config=config)
    logging.info(f"Summary saved to: {summary_path}")
    logging.info(f"Worst links saved to: {worst_path}")
    
    # Generate plots
    p1 = plot_margin_hist(df, config)
    p2 = plot_pass_or_fail(df, config)
    p3 = plot_margin_vs_length(df, config)
    p4 = plot_top_contributors(df, config)
    
    # Print paths to outputs
    print("[OK] Summary:", summary_path)
    print("[OK] Worst links:", worst_path)
    print("[OK] Plot:", p1)
    print("[OK] Plot:", p2)
    if p3:
      print("[OK] Plot:", p3)
    else:
      print("[SKIP] margin_vs_length (missing fiber_length_km — provide --input-links to merge)")
    if p4:
      print("[OK] Plot:", p4)
    else:
      print("[SKIP] top_contributors (missing top_contributor_1)")
  except Exception as e:
    raise CustomException(e, sys)

# Define main function to parse command line arguments
def main_app() -> None:
  # Argument parser for command line execution
  parser = argparse.ArgumentParser(description="Analyze Link Power Budget results from synthetic link data.")
  parser.add_argument("--results", type=Path, required=True, help="CSV output from Go engine (run/sweep)")
  parser.add_argument("--outdir", type=Path, default=Path("./exports/lpb"))
  parser.add_argument("--input-links", type=Path, default=None, help="Optional: input links CSV to merge (enables margin vs fiber_length_km plots)")
  parser.add_argument("--worst-n", type=int, default=10)
  args = parser.parse_args()
  
  config = LPBAnalysis(
    result_csv=args.results,
    output_dir=args.outdir,
    input_links_csv=args.input_links,
    worst_n=args.worst_n,
  )
  run_analysis(config=config)

if __name__ == "__main__":
  main_app()