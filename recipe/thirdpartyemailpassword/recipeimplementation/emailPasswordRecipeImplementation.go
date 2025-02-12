/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package recipeimplementation

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func MakeEmailPasswordRecipeImplementation(recipeImplementation tpepmodels.RecipeInterface) epmodels.RecipeInterface {

	signUp := func(email, password string) (epmodels.SignUpResponse, error) {
		response, err := (*recipeImplementation.SignUp)(email, password)
		if err != nil {
			return epmodels.SignUpResponse{}, err
		}
		if response.EmailAlreadyExistsError != nil {
			return epmodels.SignUpResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		}
		return epmodels.SignUpResponse{
			OK: &struct{ User epmodels.User }{
				User: epmodels.User{
					ID:         response.OK.User.ID,
					Email:      response.OK.User.Email,
					TimeJoined: response.OK.User.TimeJoined,
				},
			},
		}, nil
	}

	signIn := func(email, password string) (epmodels.SignInResponse, error) {
		response, err := (*recipeImplementation.SignIn)(email, password)
		if err != nil {
			return epmodels.SignInResponse{}, err
		}
		if response.WrongCredentialsError != nil {
			return epmodels.SignInResponse{
				WrongCredentialsError: &struct{}{},
			}, nil
		}
		return epmodels.SignInResponse{
			OK: &struct{ User epmodels.User }{
				User: epmodels.User{
					ID:         response.OK.User.ID,
					Email:      response.OK.User.Email,
					TimeJoined: response.OK.User.TimeJoined,
				},
			},
		}, nil
	}

	getUserByID := func(userId string) (*epmodels.User, error) {
		user, err := (*recipeImplementation.GetUserByID)(userId)
		if err != nil {
			return nil, err
		}
		if user == nil || user.ThirdParty != nil {
			return nil, nil
		}
		return &epmodels.User{
			ID:         user.ID,
			Email:      user.Email,
			TimeJoined: user.TimeJoined,
		}, nil
	}

	getUserByEmail := func(email string) (*epmodels.User, error) {
		users, err := (*recipeImplementation.GetUsersByEmail)(email)
		if err != nil {
			return nil, err
		}

		for _, user := range users {
			if user.ThirdParty == nil {
				return &epmodels.User{
					ID:         user.ID,
					Email:      user.Email,
					TimeJoined: user.TimeJoined,
				}, nil
			}
		}
		return nil, nil
	}

	createResetPasswordToken := func(userID string) (epmodels.CreateResetPasswordTokenResponse, error) {
		return (*recipeImplementation.CreateResetPasswordToken)(userID)
	}

	resetPasswordUsingToken := func(token, newPassword string) (epmodels.ResetPasswordUsingTokenResponse, error) {
		return (*recipeImplementation.ResetPasswordUsingToken)(token, newPassword)
	}

	updateEmailOrPassword := func(userId string, email, password *string) (epmodels.UpdateEmailOrPasswordResponse, error) {
		return (*recipeImplementation.UpdateEmailOrPassword)(userId, email, password)
	}

	return epmodels.RecipeInterface{
		SignUp:                   &signUp,
		SignIn:                   &signIn,
		GetUserByID:              &getUserByID,
		GetUserByEmail:           &getUserByEmail,
		CreateResetPasswordToken: &createResetPasswordToken,
		ResetPasswordUsingToken:  &resetPasswordUsingToken,
		UpdateEmailOrPassword:    &updateEmailOrPassword,
	}
}
