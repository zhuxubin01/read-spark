import SwiftUI

struct LoginView: View {
    let onLoginSuccess: () -> Void

    var body: some View {
        VStack(spacing: 16) {
            Text("ReadSpark")
                .font(.largeTitle)
            Button("Mock Login") {
                onLoginSuccess()
            }
            .buttonStyle(.borderedProminent)
        }
        .padding()
    }
}
