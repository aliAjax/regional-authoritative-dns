# Bug Reproduction

The worker runner executes queued jobs before their scheduled time, backoff can exceed its cap, canceled saves are persisted, and the domain runnable predicate ignores the due time.

Create a job whose next run is in the future, exercise the retry cap, and cancel a save request. The worker runs early or stores canceled work, and the worker tests fail.

