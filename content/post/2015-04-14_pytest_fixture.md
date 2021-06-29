---
title:      "Pytest Fixture"
date:       2015-04-14
categories: ["tuto"]
tags:       ["python", "dev"]
url:        "post/pytest_fixture"
---

I am a huge fan of python (one of the best language in my toolbox). And
when it comes to tests, [pytest](http://pytest.org/) is *THE* library to use.

I also use [Flask](http://flask.pocoo.org/) a lot, so today I will show you
some of my snippets.

First one the app fixture:

```python
@pytest.fixture(autouse=True)
def app():
    """Load flask in testing mode"""
    app_test = myapp
    app_test.config['TESTING'] = True
    app_test.json_encoder = my_encoder

    return app_test.test_client()
```

This create an app fixture which will be used to test the application, it
returns a test client to interact with my Flask application.

It is an adaptation of the documentation [testing skeleton]
(http://flask.pocoo.org/docs/0.10/testing/#the-testing-skeleton)

I also replace the json encoder by a custom one (it allows me to dump mock
object for example).

```python
@pytest.yield_fixture
def request_context(app):
    """Provide a Flask request context for testing purpose"""

    with app.application.test_request_context() as req_context:
        yield req_context
```

This one applies a request context on a testing function, this can be
useful if you manipulate werzeug interactions (flask request attribute mostly).
This has to be use when the
`RuntimeError: working outside of application context` error is raised.

In the most common use case you will not need the `req_context` variable during
your test. To avoid to have an unused argument, you can simply use the
`@pytest.mark.usefixtures('request_context')` decorator on your testing
function.

Last one is also an opportunity to talk about
[peewee](http://peewee.readthedocs.org/en/latest/) which is a great lightweight
ORM. Like a lot of ORM peewee provides a [transaction decorator]
(http://docs.peewee-orm.com/en/latest/peewee/transactions.html)
, in unittest you may want to test endpoints without connecting to a db.
So, here is an example on how you can replace the decorator
(or contextmanager) by a dummy one.

```python
import imp
import contextlib
import pytest

import database

@pytest.fixture(autouse=True, scope='session')
def mock_transaction():
    """Replace the atomic decorator from peewee to a noop one"""

    database.atomic = contextlib.contextmanager(lambda: (yield))
    imp.reload(module.using.db.transactions)
```

Pytest use a lot of the python flexibility and I must confess that I love it.

