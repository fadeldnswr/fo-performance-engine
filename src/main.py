'''
Docstring for main
Main application runner for fiber optic performance analysis.
'''
import os
import sys

from src.analysis.performance_analysis import run_analysis, main_app
from src.model.lpb_schema import LPBAnalysis
from src.logging.logging import logging
from src.exception.exception import CustomException

# Entry point for the application
if __name__ == "__main__":
  try:
    print("Starting Fiber Optic Performance Analysis Application")
    logging.info("Application started")
    main_app()
    logging.info("Application finished successfully")
  except Exception as e:
    raise CustomException(e, sys)
