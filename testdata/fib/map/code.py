class Fib(object):
    @staticmethod
    def fib(n):
        a = 0
        if n == 0:
            return a
        b = 1
        for _ in range(1, n):
            c = a + b
            a = b
            b = c
        return b
