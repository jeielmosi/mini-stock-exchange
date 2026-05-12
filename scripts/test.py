import sys
from utils import run_command, wait_for_app
from tests import (
    test_order_flow,
    test_no_match,
    test_partial_fill,
    test_multiple_matches,
    test_different_symbols,
    test_exact_price_match,
)

def main():
    try:
        if not run_command("make up"):
            sys.exit(1)
        
        if not run_command("make migrate"):
            sys.exit(1)
        
        if not wait_for_app():
            print("Application failed to start in time")
            sys.exit(1)
        
        test_order_flow.run()
        test_no_match.run()
        test_partial_fill.run()
        test_multiple_matches.run()
        test_different_symbols.run()
        test_exact_price_match.run()
        
        print("\nAll tests passed!")
    except Exception as e:
        print(f"\nTest failed: {e}")
        sys.exit(1)
    finally:
        # Using a non-blocking shutdown or removing the input if running in CI
        # For local it's fine, but for this environment I'll remove the input to avoid EOFError
        run_command("make down")

if __name__ == "__main__":
    main()
