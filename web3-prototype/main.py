#!/usr/bin/env python3
"""
Web3.py Learning Project - Main Entry Point

This project demonstrates Web3.py functionality through a modular, organized structure.
Run different demos to explore various Web3 concepts and features.
"""

import sys
import os

# Add src directory to Python path for imports
sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'src'))

from src.demos.basic_demo import run_basic_demo


def show_available_demos():
    """Display available demo options"""
    print("\n🚀 Web3.py Learning Project")
    print("Available demos:")
    print("  1. Basic Demo - Connections, accounts, and transactions")
    print("  2. [Coming Soon] Smart Contracts Demo")
    print("  3. [Coming Soon] Events and Logs Demo")
    print("  4. [Coming Soon] Advanced Features Demo")
    print("\nUsage: python main.py [demo_number]")
    print("       python main.py          # Run basic demo")


def main():
    """Main entry point with demo selection"""
    
    # Default to basic demo if no arguments
    demo_choice = "1"
    
    if len(sys.argv) > 1:
        demo_choice = sys.argv[1]
    
    if demo_choice == "1" or demo_choice.lower() == "basic":
        return run_basic_demo()
    elif demo_choice == "help" or demo_choice == "-h" or demo_choice == "--help":
        show_available_demos()
        return True
    else:
        print(f"❌ Unknown demo: {demo_choice}")
        show_available_demos()
        return False


if __name__ == "__main__":
    success = main()
    sys.exit(0 if success else 1)