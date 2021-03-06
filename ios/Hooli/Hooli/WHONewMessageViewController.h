//
//  WHONewMessageViewController.h
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol WHOMessageProtocol <NSObject>
- (void)receivedNewMessage:(NSString* )message;
@end

@interface WHONewMessageViewController : UIViewController <UITextViewDelegate>

- (IBAction)messageView:(id)sender;
@property (nonatomic, strong) id<WHOMessageProtocol> delegate;
@property (strong, nonatomic) IBOutlet UITextView *messageField;
@property (strong, nonatomic) IBOutlet UILabel *placeholderLabel;

@end
