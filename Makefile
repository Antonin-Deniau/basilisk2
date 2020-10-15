.PHONY: test

test:
	./runtest.py --debug=test tests/step2_eval.mal ./basilisk/repl.py
