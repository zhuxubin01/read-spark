import SwiftUI

struct ContentView: View {
    @State private var isLoggedIn = KeychainManager.shared.accessToken != nil

    var body: some View {
        if isLoggedIn {
            AppTabView()
        } else {
            LoginView(onLoginSuccess: { isLoggedIn = true })
        }
    }
}
