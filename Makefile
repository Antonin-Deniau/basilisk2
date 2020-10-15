.PHONY: test

test1:
	./runtest.py --debug=test tests/step2_eval.mal ./basilisk/repl.py

test2:
	./runtest.py --debug=test tests/step3_env.mal ./basilisk/repl.py

test3:
	./runtest.py --debug=test tests/step4_if_fn_do.mal ./basilisk/repl.py
