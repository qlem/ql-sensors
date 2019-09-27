##
## Created by Clément Dommerc 11/04/2019
##

NAME		:=	qlsensors
CC		:=	go
CP		:=	cp
RM		:=	rm -f

SRC_DIR		:=	src
INSTALL_DIR	:=	/usr/local/bin

SRC		:=	main.go		\
			ncurses.go	\
			print.go	\
			sensors.go	\
			cpu.go
SRC	:=	$(addprefix $(SRC_DIR)/, $(SRC))

all:
	$(CC) build -o bin/$(NAME) $(SRC)

install:
	$(CP) bin/$(NAME) $(INSTALL_DIR)

uninstall:
	$(RM) $(INSTALL_DIR)/$(NAME)

clean:
	$(RM) bin/$(NAME)
