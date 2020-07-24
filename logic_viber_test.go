package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWeHaveNewMessageAfterUnsubscribe(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()
	userId := "124"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, "Добрый день, Vasya. Добро пожаловать. Вы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"1. Да", "2. Нет"}, reply.options)

	text := newTextCallback(t, userId, "1. Да")
	require.Equal(t, text.User.Id, userId)
	reply, err = generateReplyFor(p, s, text)
	require.NoError(t, err)
	require.Equal(t, "Укажите, пожалуйста, Ваш возраст", reply.text)

	reply, err = generateReplyFor(p, s, newUnsubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, "Добрый день, Vasya. Добро пожаловать. Укажите, пожалуйста, Ваш возраст", reply.text)
	require.Equal(t, []string{"1. меньше 18", "2. 18-24", "3. 25-34", "4. 35-44", "5. 45-54", "6. 55+"}, reply.options)
}

func TestUserFlowCaseSensitive(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	userId := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, "Добрый день, Vasya. Добро пожаловать. Вы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"1. Да", "2. Нет"}, reply.options)

	text := newTextCallback(t, userId, "Привет")
	require.Equal(t, text.User.Id, userId)
	reply, err = generateReplyFor(p, s, text)
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста выберите предложенный ответ. Вы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"1. Да", "2. Нет"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. Да"))
	require.NoError(t, err)
	require.Equal(t, "Укажите, пожалуйста, Ваш возраст", reply.text)
	require.Equal(t, []string{"1. меньше 18", "2. 18-24", "3. 25-34", "4. 35-44", "5. 45-54", "6. 55+"}, reply.options)

	user, err := s.fromPersisted(userId)
	require.NoError(t, err)
	require.Equal(t, user.Id, userId)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Clear"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Your storage cleared with <nil>")

	user, err = s.fromPersisted(userId)
	require.NoError(t, err)
	require.Nil(t, user)
}

func TestCaseInsensitive(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	caseSensitive = false
	userId := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Добрый день, Vasya. Добро пожаловать. Вы гражданин Республики Беларусь?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста выберите предложенный ответ. Вы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "2. нет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Только граждание Беларуси могут принимать участие! Вы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. да"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. меньше 18"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Вам должно быть 18 или больше. Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "4. 35-44"))
	require.NoError(t, err)
	require.Equal(t, "Примете ли Вы участие в предстоящих выборах Президента?", reply.text)
	require.Equal(t, []string{"1. Да", "2. Нет", "3. Затрудняюсь ответить"}, reply.options)
	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. Да"))
	require.NoError(t, err)
	require.Equal(t, "За кого Вы планируете проголосовать?", reply.text)
	require.Equal(t, []string{"1. Дмитриев", "2. Канопацкая", "3. Лукашенко", "4. Тихановская", "5. Черечень", "6. Против всех", "7. Затрудняюсь ответить"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. Дмитриев"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш пол")
	require.Equal(t, reply.options, []string{"1. Мужской", "2. Женский"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. мужской"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Выберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск")
	require.Equal(t, reply.options, []string{"1. Брестская", "2. Витебская", "3. Гомельская", "4. Гродненская", "5. Минская", "6. Могилевская", "7. Минск", "8. Проживаю за пределами РБ"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "2. витебская"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "К какому типу относится ваш населенный пункт? Если живете в Минске - выберите Минск.")
	require.Equal(t, reply.options, []string{"1. Областной центр", "2. Город или городской поселок", "3. Агрогородок / Село / Деревня", "4. Проживаю за пределами РБ"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "3. Агрогородок / село / Деревня"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Ваш уровень образования?")
	require.Equal(t, reply.options, []string{"1. Базовое / Среднее общее (школа)", "2. Профессионально-техническое", "3. Среднее специальное", "4. Высшее", "6. Другое"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "3. Среднее специальное"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, общий совокупный доход вашей семьи (включая пенсии, стипендии, надбавки и прочее)")
	require.Equal(t, reply.options, []string{"1. До 500 бел. руб.", "2. 500 - 1000 бел. руб.", "3. 1000 - 2000 бел. руб.", "4. Выше 2000 бел.руб.", "5. Не хочу отвечать на этот вопрос"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. До 500 бел. руб."))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Спасибо за голосование! Уже проголосовало 1 человек")

	user, err := s.fromPersisted(userId)
	require.NoError(t, err)
	require.Equal(t, user.Id, userId)
	require.Equal(t, user.Properties["age"], "4. 35-44")
	require.Equal(t, user.Level, 9)
	require.Equal(t, user.Candidate, "1. Дмитриев")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Спасибо за голосование! Уже проголосовало 1 человек")

	reply, err = generateReplyFor(p, s, newUnsubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, "Добрый день, Vasya. Добро пожаловать. Спасибо за голосование! Уже проголосовало 1 человек", reply.text)
}

func TestUserFlow(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	userId := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Добрый день, Vasya. Добро пожаловать. Вы гражданин Республики Беларусь?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста выберите предложенный ответ. Вы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "2. Нет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Только граждание Беларуси могут принимать участие! Вы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. Да"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. меньше 18"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Вам должно быть 18 или больше. Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "4. 35-44"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Примете ли Вы участие в предстоящих выборах Президента?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет", "3. Затрудняюсь ответить"})

	user, err := s.fromPersisted(userId)
	require.NoError(t, err)

	require.Equal(t, user.Id, userId)
	fmt.Println(user)
	require.Equal(t, user.Properties["age"], "4. 35-44")
	require.Equal(t, user.Level, 2)

	seenCallback := newSeenCallback(t, userId)
	require.Equal(t, seenCallback.User.Id, userId)
	reply, err = generateReplyFor(p, s, seenCallback)
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. Да"))
	require.NoError(t, err)
	require.Equal(t, "За кого Вы планируете проголосовать?", reply.text)
	require.Equal(t, []string{"1. Дмитриев", "2. Канопацкая", "3. Лукашенко", "4. Тихановская", "5. Черечень", "6. Против всех", "7. Затрудняюсь ответить"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "3. Лукашенко"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш пол")
	require.Equal(t, reply.options, []string{"1. Мужской", "2. Женский"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1. Мужской"))
	require.NoError(t, err)

	user, err = s.fromPersisted(userId)
	require.NoError(t, err)

	require.Equal(t, user.Id, userId)
	require.Equal(t, user.Properties["age"], "4. 35-44")
	require.Equal(t, 5, user.Level)
	require.Equal(t, "3. Лукашенко", user.Candidate)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста выберите предложенный ответ. Выберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск", reply.text)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста выберите предложенный ответ. Выберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск", reply.text)

	subscribe := newSubscribeCallback(t, userId)
	user, err = s.Obtain(userId)
	require.NoError(t, err)
	reply, err = generateReplyFor(p, s, subscribe)
	require.NoError(t, err)
	require.Empty(t, reply)

	reply, err = generateReplyFor(p, s, newSeenCallback(t, userId))
	require.NoError(t, err)
	require.Empty(t, reply)
}

func TestNumericReplies(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	userId := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Добрый день, Vasya. Добро пожаловать. Вы гражданин Республики Беларусь?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста выберите предложенный ответ. Вы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "2"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Только граждание Беларуси могут принимать участие! Вы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "120"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста выберите предложенный ответ. Вы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "1"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "4"))
	require.NoError(t, err)
	require.Equal(t, "Примете ли Вы участие в предстоящих выборах Президента?", reply.text)

	user, err := s.fromPersisted(userId)
	require.NoError(t, err)

	require.Equal(t, user.Id, userId)
	require.Equal(t, 2, user.Level)
	require.Equal(t, user.Properties["age"], "4. 35-44")
}

func newUnsubscribeCallback(t *testing.T, id string) *ViberCallback {
	json := `{"event":"unsubscribed","timestamp":1595347885535,"chat_hostname":"SN-376_","user_id":"%s","message_token":5466394919049723652}`
	validJson := fmt.Sprintf(json, id)

	ret, err := parseCallback([]byte(validJson))
	require.NoError(t, err)

	return ret
}

func newSubscribeCallback(t *testing.T, id string) *ViberCallback {
	c := &ViberCallback{
		Event: "subscribed",
		User: User{
			Id:   id,
			Name: "Vasya",
		},
	}

	b, err := json.Marshal(c)
	require.NoError(t, err)

	ret, err := parseCallback(b)
	require.NoError(t, err)

	return ret
}

func newTextCallback(t *testing.T, id string, text string) *ViberCallback {
	json := `{"event":"message","sender":{"id":"%s","Name":"Vasya"},"message":{"type":"text","text":"%s"}}`

	validJson := fmt.Sprintf(json, id, text)

	ret, err := parseCallback([]byte(validJson))
	require.NoError(t, err)

	return ret
}

func newSeenCallback(t *testing.T, id string) *ViberCallback {
	json := `{"event":"seen","user_id":"%s"}`

	validJson := fmt.Sprintf(json, id)

	ret, err := parseCallback([]byte(validJson))
	require.NoError(t, err)

	return ret
}
