NAME = collectator-armv7l
BUILDER = $(NAME)-builder

all: $(NAME)

$(NAME): Dockerfile collectator.go
	docker rm $(BUILDER) 2>/dev/null || true
	docker build -t $(BUILDER) .
	docker run --name $(BUILDER) $(BUILDER) echo hej
	docker cp $(BUILDER):/app/collectator $(NAME)
	docker rm $(BUILDER)

fclean:
	rm -f $(NAME)
