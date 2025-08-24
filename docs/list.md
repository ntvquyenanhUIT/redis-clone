# Redis List Command Documentation

This document provides detailed information about the supported Redis list commands in this project.

Our server implements lists using a doubly linked list, which ensures that adding or removing elements from either the head (left) or tail (right) of the list is a very fast, O(1) operation.

---

## Supported Commands

### `LPUSH`
Inserts one or more elements at the head (left side) of the list.

-   **Syntax**: `LPUSH key element [element ...]`
-   **Arguments**:
    -   `key`: The key holding the list.
    -   `element`: One or more elements to prepend to the list.
-   **Returns**: An integer representing the new length of the list after the push operation.
-   **Example**:
    ```bash
    127.0.0.1:6379> LPUSH mylist "world"
    (integer) 1
    127.0.0.1:6379> LPUSH mylist "hello" "and"
    (integer) 3
    ```

### `RPUSH`
Inserts one or more elements at the tail (right side) of the list.

-   **Syntax**: `RPUSH key element [element ...]`
-   **Arguments**:
    -   `key`: The key holding the list.
    -   `element`: One or more elements to append to the list.
-   **Returns**: An integer representing the new length of the list after the push operation.
-   **Example**:
    ```bash
    127.0.0.1:6379> RPUSH mylist "one"
    (integer) 1
    127.0.0.1:6379> RPUSH mylist "two" "three"
    (integer) 3
    ```

### `LPOP`
Removes and returns one or more elements from the head (left side) of the list.

-   **Syntax**: `LPOP key [count]`
-   **Arguments**:
    -   `key`: The key holding the list.
    -   `count` (optional): The number of elements to pop. Defaults to 1.
-   **Returns**:
    -   Without `count`: A bulk string reply with the value of the popped element, or `(nil)` if the list is empty.
    -   With `count`: An array of bulk strings with the popped elements.
-   **Example**:
    ```bash
    127.0.0.1:6379> RPUSH mylist "a" "b" "c"
    (integer) 3
    127.0.0.1:6379> LPOP mylist
    "a"
    127.0.0.1:6379> LPOP mylist 2
    1) "b"
    2) "c"
    ```

### `BLPOP`
A blocking version of `LPOP`. It blocks the connection when no elements are available in the list, waiting for an element to be pushed or until the timeout is reached.

-   **Syntax**: `BLPOP key timeout`
-   **Arguments**:
    -   `key`: The key to watch.
    -   `timeout`: A floating-point number specifying the maximum number of seconds to block. A timeout of `0` will block indefinitely.
-   **Returns**: An array reply containing the key name and the popped element, or `(nil)` if the timeout is reached.
-   **Example**:
    ```bash
    # Terminal 1
    127.0.0.1:6379> BLPOP mylist 5
    # (blocks)

    # Terminal 2 (within 5 seconds)
    127.0.0.1:6379> LPUSH mylist "new-item"
    (integer) 1

    # Terminal 1 (unblocks and returns)
    1) "mylist"
    2) "new-item"
    ```

### `LLEN`
Returns the length of the list stored at `key`.

-   **Syntax**: `LLEN key`
-   **Returns**: An integer reply with the length of the list, or `0` if the key does not exist.
-   **Example**:
    ```bash
    127.0.0.1:6379> LPUSH mylist "a" "b"
    (integer) 2
    127.0.0.1:6379> LLEN mylist
    (integer) 2
    ```

### `LRANGE`
Returns a range of elements from the list. The range is specified by `start` and `end` offsets, which are zero-based. These offsets can also be negative numbers to indicate positions from the end of the list.

-   **Syntax**: `LRANGE key start end`
-   **Returns**: An array of elements within the specified range.
-   **Example**:
    ```bash
    127.0.0.1:6379> RPUSH mylist "one" "two" "three" "four"
    (integer) 4
    127.0.0.1:6379> LRANGE mylist 0 2
    1) "one"
    2) "two"
    3) "three"
    127.0.0.1:6379> LRANGE mylist -2 -1
    1) "three"
    2) "four"
    ```