import SwiftUI

struct LoginView: View {
    @StateObject private var viewModel = LoginViewModel()
    let onLoginSuccess: () -> Void

    var body: some View {
        VStack(spacing: 20) {
            Text("ReadSpark")
                .font(.largeTitle)
                .fontWeight(.bold)

            TextField("手机号", text: $viewModel.phone)
                .textFieldStyle(RoundedBorderTextFieldStyle())
                .keyboardType(.phonePad)

            SecureField("验证码 (测试用: 123456)", text: $viewModel.code)
                .textFieldStyle(RoundedBorderTextFieldStyle())

            if let error = viewModel.error {
                Text(error)
                    .foregroundColor(.red)
                    .font(.caption)
            }

            Button(action: { viewModel.login() }) {
                if viewModel.isLoading {
                    ProgressView()
                } else {
                    Text("登录 / 注册")
                        .frame(maxWidth: .infinity)
                }
            }
            .buttonStyle(.borderedProminent)
            .disabled(viewModel.isLoading)

            Spacer()
        }
        .padding()
        .onChange(of: viewModel.isLoggedIn) { loggedIn in
            if loggedIn {
                onLoginSuccess()
            }
        }
    }
}
